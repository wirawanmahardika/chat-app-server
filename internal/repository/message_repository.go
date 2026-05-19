package repository

import (
	"chatapp/database/model"
	"context"
	"time"

	"gorm.io/gorm"
)

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db}
}

type MessageRepository struct {
	db *gorm.DB
}

func (r *MessageRepository) CreateMessage(c context.Context, conversationID, senderID, receiverID, text string) (*model.Message, error) {
	message := model.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		ReceiverID:     receiverID,
		Text:           text,
		Read:           false,
		Delivered:      true, // since HTTP is instant delivery for this REST fallback
	}

	err := r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&message).Error; err != nil {
			return err
		}

		// Update Conversation
		updates := map[string]interface{}{
			"last_message_id": message.ID,
			"updated_at":      time.Now(),
		}
		if err := tx.Model(&model.Conversation{}).Where("id = ?", conversationID).Updates(updates).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *MessageRepository) GetMessages(c context.Context, conversationID string, limit, offset int) ([]model.Message, int64, error) {
	var messages []model.Message
	var total int64

	// Count total messages
	if err := r.db.WithContext(c).Model(&model.Message{}).Where("conversation_id = ?", conversationID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch messages ordered by createdAt DESC for cursor/offset pagination
	err := r.db.WithContext(c).
		Where("conversation_id = ?", conversationID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error

	return messages, total, err
}

func (r *MessageRepository) MarkAsRead(c context.Context, messageID string, receiverID string) error {
	return r.db.WithContext(c).
		Model(&model.Message{}).
		Where("id = ? AND receiver_id = ?", messageID, receiverID).
		Update("read", true).Error
}

func (r *MessageRepository) GetByID(c context.Context, id string) (*model.Message, error) {
	var message model.Message
	err := r.db.WithContext(c).First(&message, "id = ?", id).Error
	return &message, err
}

