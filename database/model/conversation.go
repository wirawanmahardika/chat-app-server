package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Conversation struct {
	ID            string    `gorm:"primaryKey;type:text" json:"id"`
	Type          string    `gorm:"default:private" json:"type"`
	LastMessageID *string   `json:"lastMessageId"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	Participants []ConversationParticipant `gorm:"foreignKey:ConversationID" json:"-"`
	Messages     []Message                 `gorm:"foreignKey:ConversationID" json:"-"`
	LastMessage  *Message                  `gorm:"foreignKey:LastMessageID" json:"-"`
}

func (c *Conversation) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	return nil
}
