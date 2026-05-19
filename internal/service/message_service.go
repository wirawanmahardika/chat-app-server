package service

import (
	"chatapp/database/model"
	"chatapp/internal/dto"
	"chatapp/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func NewMessageService(
	messageRepo *repository.MessageRepository,
	convoRepo *repository.ConversationRepository,
	userRepo *repository.UserRepository,
	friendshipRepo *repository.FriendshipRepository,
) *MessageService {
	return &MessageService{messageRepo, convoRepo, userRepo, friendshipRepo}
}

type MessageService struct {
	messageRepository      *repository.MessageRepository
	conversationRepository *repository.ConversationRepository
	userRepository          *repository.UserRepository
	friendshipRepository    *repository.FriendshipRepository
}

func (s *MessageService) GetConversations(c context.Context, userID string) ([]dto.ConversationResponse, error) {
	convos, err := s.conversationRepository.GetConversations(c, userID)
	if err != nil {
		return nil, err
	}

	var res []dto.ConversationResponse
	for _, convo := range convos {
		// Find the other participant (friend)
		var friend *model.User
		for _, part := range convo.Participants {
			if part.UserID != userID {
				friend = part.User
				break
			}
		}

		// If no other participant (e.g. self chat or corrupted data), skip
		if friend == nil {
			continue
		}

		var lastMsg *string
		var lastMsgTime *string

		if convo.LastMessage != nil {
			lastMsg = &convo.LastMessage.Text
			tStr := convo.LastMessage.CreatedAt.Format(time.RFC3339)
			lastMsgTime = &tStr
		}

		res = append(res, dto.ConversationResponse{
			ID:              friend.ID,
			Name:            friend.Name,
			Avatar:          friend.Avatar,
			LastMessage:     lastMsg,
			Online:          friend.Online,
			LastSeen:        friend.LastSeen.Format(time.RFC3339),
			LastMessageTime: lastMsgTime,
		})
	}

	if res == nil {
		res = []dto.ConversationResponse{}
	}
	return res, nil
}

func (s *MessageService) GetMessages(c context.Context, userID, conversationID string, limit, offset int) (*dto.MessagesWithPagination, error) {
	// Verify user is participant
	isPart, err := s.conversationRepository.IsParticipant(c, conversationID, userID)
	if err != nil {
		return nil, err
	}
	if !isPart {
		return nil, fiber.NewError(fiber.StatusForbidden, "You are not a participant of this conversation")
	}

	messages, total, err := s.messageRepository.GetMessages(c, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}

	var resMessages []dto.MessageResponse
	for _, msg := range messages {
		resMessages = append(resMessages, dto.MessageResponse{
			ID:             msg.ID,
			ConversationID: msg.ConversationID,
			SenderID:       msg.SenderID,
			ReceiverID:     msg.ReceiverID,
			Text:           msg.Text,
			Read:           msg.Read,
			Delivered:      msg.Delivered,
			CreatedAt:      msg.CreatedAt.Format(time.RFC3339),
		})
	}

	if resMessages == nil {
		resMessages = []dto.MessageResponse{}
	}

	return &dto.MessagesWithPagination{
		Success: true,
		Data:    resMessages,
		Pagination: dto.PaginationInfo{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
	}, nil
}

func (s *MessageService) SendMessage(c context.Context, senderID, receiverID, text string) (*dto.MessageResponse, error) {
	if senderID == receiverID {
		return nil, fiber.NewError(fiber.StatusBadRequest, "You cannot send a message to yourself")
	}

	// Verify receiver exists
	_, err := s.userRepository.GetByID(c, receiverID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Receiver not found")
		}
		return nil, err
	}

	// Verify friendship exists (must be friends)
	friendship, err := s.friendshipRepository.GetFriendshipBetween(c, senderID, receiverID)
	if err != nil || friendship.Status != "accepted" {
		return nil, fiber.NewError(fiber.StatusForbidden, "You can only send messages to friends")
	}

	// Get or Create Conversation
	convo, err := s.conversationRepository.GetOrCreatePrivateConversation(c, senderID, receiverID)
	if err != nil {
		return nil, err
	}

	// Save message
	msg, err := s.messageRepository.CreateMessage(c, convo.ID, senderID, receiverID, text)
	if err != nil {
		return nil, err
	}

	return &dto.MessageResponse{
		ID:             msg.ID,
		ConversationID: msg.ConversationID,
		SenderID:       msg.SenderID,
		ReceiverID:     msg.ReceiverID,
		Text:           msg.Text,
		Read:           msg.Read,
		Delivered:      msg.Delivered,
		CreatedAt:      msg.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *MessageService) MarkAsRead(c context.Context, userID, messageID string) error {
	// Verify message exists and userID is the receiver
	msg, err := s.messageRepository.GetByID(c, messageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Message not found")
		}
		return err
	}

	if msg.ReceiverID != userID {
		return fiber.NewError(fiber.StatusForbidden, "You cannot mark this message as read")
	}

	return s.messageRepository.MarkAsRead(c, messageID, userID)
}
