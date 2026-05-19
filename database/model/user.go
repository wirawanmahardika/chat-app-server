package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"primaryKey;type:text" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique;not null;index" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Avatar    *string   `json:"avatar"`
	Online    bool      `gorm:"default:false" json:"online"`
	LastSeen  time.Time `json:"lastSeen"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	Participants []ConversationParticipant `gorm:"foreignKey:UserID" json:"-"`
	Messages     []Message                 `gorm:"foreignKey:SenderID" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return err
	}
	u.Password = string(hashPassword)

	return nil
}

func (u *User) ValidatePassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return false
	}
	return true
}
