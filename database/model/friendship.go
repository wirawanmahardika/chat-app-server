package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Friendship struct {
	ID         string    `gorm:"primaryKey;type:text" json:"id"`
	SenderID   string    `gorm:"not null;index" json:"fromId"`
	ReceiverID string    `gorm:"not null;index" json:"toId"`
	Status     string    `gorm:"default:pending" json:"status"` // pending, accepted, declined
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	Sender   User `gorm:"foreignKey:SenderID" json:"-"`
	Receiver User `gorm:"foreignKey:ReceiverID" json:"-"`
}

func (f *Friendship) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.NewString()
	}
	return nil
}
