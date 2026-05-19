package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConversationParticipant struct {
	ID             string `gorm:"primaryKey;type:text" json:"id"`
	ConversationID string `gorm:"not null;index" json:"conversationId"`
	UserID         string `gorm:"not null;index" json:"userId"`

	User         *User         `gorm:"foreignKey:UserID" json:"-"`
	Conversation *Conversation `gorm:"foreignKey:ConversationID" json:"-"`
}

func (cp *ConversationParticipant) BeforeCreate(tx *gorm.DB) error {
	if cp.ID == "" {
		cp.ID = uuid.NewString()
	}
	return nil
}
