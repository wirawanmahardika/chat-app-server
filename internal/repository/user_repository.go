package repository

import (
	"chatapp/database/model"
	"context"
	"time"

	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) Create(c context.Context, name, email, password string) (*model.User, error) {
	user := model.User{
		Name:     name,
		Email:    email,
		Password: password,
		Online:   false,
		LastSeen: time.Now(),
	}

	result := r.db.WithContext(c).Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (r *UserRepository) CountByEmail(c context.Context, email string) (int64, error) {
	var count int64
	if err := r.db.WithContext(c).Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserRepository) GetByID(c context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(c).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(c context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(c).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Search(c context.Context, q string, limit int, excludeUserID string) ([]model.User, error) {
	var users []model.User
	query := r.db.WithContext(c).
		Model(&model.User{}).
		Where("id != ?", excludeUserID).
		Where("name LIKE ? OR email LIKE ?", "%"+q+"%", "%"+q+"%").
		Limit(limit)

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) UpdateOnlineStatus(c context.Context, userID string, online bool) error {
	updates := map[string]interface{}{
		"online":    online,
		"last_seen": time.Now(),
	}
	return r.db.WithContext(c).Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error;
}
