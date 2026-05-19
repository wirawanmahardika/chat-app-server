package repository

import (
	"chatapp/database/model"
	"context"

	"gorm.io/gorm"
)

func NewConversationRepository(db *gorm.DB) *ConversationRepository {
	return &ConversationRepository{db}
}

type ConversationRepository struct {
	db *gorm.DB
}

func (r *ConversationRepository) GetOrCreatePrivateConversation(c context.Context, user1ID, user2ID string) (*model.Conversation, error) {
	var cp model.ConversationParticipant
	// Search if private conversation already exists
	err := r.db.WithContext(c).Raw(`
		SELECT cp1.* FROM conversation_participants cp1
		JOIN conversation_participants cp2 ON cp1.conversation_id = cp2.conversation_id
		JOIN conversations c ON cp1.conversation_id = c.id
		WHERE cp1.user_id = ? AND cp2.user_id = ? AND c.type = 'private'
		LIMIT 1
	`, user1ID, user2ID).Scan(&cp).Error

	if err == nil && cp.ConversationID != "" {
		var conversation model.Conversation
		if err := r.db.WithContext(c).First(&conversation, "id = ?", cp.ConversationID).Error; err == nil {
			return &conversation, nil
		}
	}

	// Create new conversation
	var conversation model.Conversation
	conversation.Type = "private"

	err = r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&conversation).Error; err != nil {
			return err
		}

		p1 := model.ConversationParticipant{
			ConversationID: conversation.ID,
			UserID:         user1ID,
		}
		p2 := model.ConversationParticipant{
			ConversationID: conversation.ID,
			UserID:         user2ID,
		}

		if err := tx.Create(&p1).Error; err != nil {
			return err
		}
		if err := tx.Create(&p2).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (r *ConversationRepository) GetConversations(c context.Context, userID string) ([]model.Conversation, error) {
	var participants []model.ConversationParticipant
	// Find all conversation IDs where the user is a participant
	err := r.db.WithContext(c).
		Where("user_id = ?", userID).
		Find(&participants).Error

	if err != nil {
		return nil, err
	}

	if len(participants) == 0 {
		return []model.Conversation{}, nil
	}

	var conversationIDs []string
	for _, p := range participants {
		conversationIDs = append(conversationIDs, p.ConversationID)
	}

	var conversations []model.Conversation
	err = r.db.WithContext(c).
		Preload("LastMessage").
		Preload("Participants.User").
		Where("id IN ?", conversationIDs).
		Order("updated_at DESC").
		Find(&conversations).Error

	return conversations, err
}

func (r *ConversationRepository) IsParticipant(c context.Context, conversationID, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(c).
		Model(&model.ConversationParticipant{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Count(&count).Error

	return count > 0, err
}
