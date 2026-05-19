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

func NewFriendService(
	friendshipRepo *repository.FriendshipRepository,
	userRepo *repository.UserRepository,
	convoRepo *repository.ConversationRepository,
) *FriendService {
	return &FriendService{friendshipRepo, userRepo, convoRepo}
}

type FriendService struct {
	friendshipRepository *repository.FriendshipRepository
	userRepository       *repository.UserRepository
	conversationRepository *repository.ConversationRepository
}

func (s *FriendService) SendRequest(c context.Context, senderID, friendID string) error {
	if senderID == friendID {
		return fiber.NewError(fiber.StatusBadRequest, "You cannot send a friend request to yourself")
	}

	// Verify friend exists
	_, err := s.userRepository.GetByID(c, friendID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return err
	}

	// Check existing friendship
	existing, err := s.friendshipRepository.GetFriendshipBetween(c, senderID, friendID)
	if err == nil {
		if existing.Status == "accepted" {
			return fiber.NewError(fiber.StatusConflict, "You are already friends")
		}
		if existing.Status == "pending" {
			if existing.SenderID == senderID {
				return fiber.NewError(fiber.StatusConflict, "Friend request already sent")
			} else {
				// The other user already sent a request, so accept it!
				return s.AcceptRequestAuto(c, senderID, existing.ID)
			}
		}
		// If declined, reset to pending
		if existing.Status == "declined" {
			err = s.friendshipRepository.UpdateStatus(c, existing.ID, "pending")
			if err != nil {
				return err
			}
			// Swap sender and receiver to match the new request
			existing.SenderID = senderID
			existing.ReceiverID = friendID
			return s.friendshipRepository.Save(c, existing)
		}
	}

	_, err = s.friendshipRepository.CreateRequest(c, senderID, friendID)
	return err
}

func (s *FriendService) AcceptRequestAuto(c context.Context, userID, requestID string) error {
	_, err := s.AcceptRequest(c, userID, requestID)
	return err
}

func (s *FriendService) GetPendingRequests(c context.Context, userID string) ([]dto.FriendRequestResponse, error) {
	requests, err := s.friendshipRepository.GetPendingRequests(c, userID)
	if err != nil {
		return nil, err
	}

	var res []dto.FriendRequestResponse
	for _, req := range requests {
		var avatar *string
		if req.Sender.Avatar != nil {
			avatar = req.Sender.Avatar
		}
		res = append(res, dto.FriendRequestResponse{
			ID:         req.ID,
			FromID:     req.SenderID,
			FromName:   req.Sender.Name,
			FromAvatar: avatar,
			ToID:       req.ReceiverID,
			Status:     req.Status,
			CreatedAt:  req.CreatedAt.Format(time.RFC3339),
		})
	}

	if res == nil {
		res = []dto.FriendRequestResponse{}
	}
	return res, nil
}

func (s *FriendService) AcceptRequest(c context.Context, userID, requestID string) (*dto.FriendResponse, error) {
	friendship, err := s.friendshipRepository.GetFriendship(c, requestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Friend request not found")
		}
		return nil, err
	}

	if friendship.ReceiverID != userID {
		return nil, fiber.NewError(fiber.StatusForbidden, "You cannot accept this friend request")
	}

	if friendship.Status != "pending" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Friend request is not pending")
	}

	// Update status
	err = s.friendshipRepository.UpdateStatus(c, requestID, "accepted")
	if err != nil {
		return nil, err
	}

	// Get or Create Private Conversation
	_, err = s.conversationRepository.GetOrCreatePrivateConversation(c, friendship.SenderID, friendship.ReceiverID)
	if err != nil {
		return nil, err
	}

	// Return the friend details (the sender of the request)
	friend, err := s.userRepository.GetByID(c, friendship.SenderID)
	if err != nil {
		return nil, err
	}

	return &dto.FriendResponse{
		ID:       friend.ID,
		Name:     friend.Name,
		Avatar:   friend.Avatar,
		Online:   friend.Online,
		LastSeen: friend.LastSeen.Format(time.RFC3339),
	}, nil
}

func (s *FriendService) DeclineRequest(c context.Context, userID, requestID string) error {
	friendship, err := s.friendshipRepository.GetFriendship(c, requestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Friend request not found")
		}
		return err
	}

	if friendship.ReceiverID != userID {
		return fiber.NewError(fiber.StatusForbidden, "You cannot decline this friend request")
	}

	if friendship.Status != "pending" {
		return fiber.NewError(fiber.StatusBadRequest, "Friend request is not pending")
	}

	return s.friendshipRepository.UpdateStatus(c, requestID, "declined")
}

func (s *FriendService) GetFriends(c context.Context, userID string) ([]dto.FriendResponse, error) {
	friendships, err := s.friendshipRepository.GetFriends(c, userID)
	if err != nil {
		return nil, err
	}

	var res []dto.FriendResponse
	for _, fs := range friendships {
		var friend model.User
		if fs.SenderID == userID {
			friend = fs.Receiver
		} else {
			friend = fs.Sender
		}

		// Find last message if any conversation exists
		var lastMsgText *string
		convo, err := s.conversationRepository.GetOrCreatePrivateConversation(c, userID, friend.ID)
		if err == nil && convo.LastMessageID != nil {
			if msg, err := s.friendshipRepository.GetMessageByID(c, *convo.LastMessageID); err == nil {
				lastMsgText = &msg.Text
			}
		}

		res = append(res, dto.FriendResponse{
			ID:          friend.ID,
			Name:        friend.Name,
			Avatar:      friend.Avatar,
			LastMessage: lastMsgText,
			Online:      friend.Online,
			LastSeen:    friend.LastSeen.Format(time.RFC3339),
		})
	}

	if res == nil {
		res = []dto.FriendResponse{}
	}
	return res, nil
}

func (s *FriendService) RemoveFriend(c context.Context, userID, friendID string) error {
	friendship, err := s.friendshipRepository.GetFriendshipBetween(c, userID, friendID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Friend connection not found")
		}
		return err
	}

	return s.friendshipRepository.DeleteFriendship(c, friendship.ID)
}
