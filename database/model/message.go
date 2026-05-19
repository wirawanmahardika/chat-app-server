package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID             string    `gorm:"primaryKey;type:text" json:"id"`
	ConversationID string    `gorm:"not null;index" json:"conversationId"`
	SenderID       string    `gorm:"not null;index" json:"senderId"`
	ReceiverID     string    `gorm:"not null;index" json:"receiverId"`
	Text           string    `gorm:"not null;type:text" json:"text"`
	Read           bool      `gorm:"default:false" json:"read"`
	Delivered      bool      `gorm:"default:false" json:"delivered"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`

	Sender   User `gorm:"foreignKey:SenderID" json:"-"`
	Receiver User `gorm:"foreignKey:ReceiverID" json:"-"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return nil
}
