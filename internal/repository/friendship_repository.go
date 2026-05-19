package repository

import (
	"chatapp/database/model"
	"context"

	"gorm.io/gorm"
)

func NewFriendshipRepository(db *gorm.DB) *FriendshipRepository {
	return &FriendshipRepository{db}
}

type FriendshipRepository struct {
	db *gorm.DB
}

func (r *FriendshipRepository) CreateRequest(c context.Context, senderID, receiverID string) (*model.Friendship, error) {
	friendship := model.Friendship{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     "pending",
	}

	if err := r.db.WithContext(c).Create(&friendship).Error; err != nil {
		return nil, err
	}
	return &friendship, nil
}

func (r *FriendshipRepository) GetPendingRequests(c context.Context, userID string) ([]model.Friendship, error) {
	var requests []model.Friendship
	err := r.db.WithContext(c).
		Preload("Sender").
		Where("receiver_id = ? AND status = ?", userID, "pending").
		Find(&requests).Error

	return requests, err
}

func (r *FriendshipRepository) GetFriendship(c context.Context, id string) (*model.Friendship, error) {
	var friendship model.Friendship
	err := r.db.WithContext(c).
		Preload("Sender").
		Preload("Receiver").
		First(&friendship, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return &friendship, nil
}

func (r *FriendshipRepository) GetFriendshipBetween(c context.Context, user1ID, user2ID string) (*model.Friendship, error) {
	var friendship model.Friendship
	err := r.db.WithContext(c).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", user1ID, user2ID, user2ID, user1ID).
		First(&friendship).Error

	if err != nil {
		return nil, err
	}
	return &friendship, nil
}

func (r *FriendshipRepository) GetFriends(c context.Context, userID string) ([]model.Friendship, error) {
	var friendships []model.Friendship
	err := r.db.WithContext(c).
		Preload("Sender").
		Preload("Receiver").
		Where("status = ? AND (sender_id = ? OR receiver_id = ?)", "accepted", userID, userID).
		Find(&friendships).Error

	return friendships, err
}

func (r *FriendshipRepository) UpdateStatus(c context.Context, requestID string, status string) error {
	return r.db.WithContext(c).
		Model(&model.Friendship{}).
		Where("id = ?", requestID).
		Update("status", status).Error
}

func (r *FriendshipRepository) DeleteFriendship(c context.Context, friendshipID string) error {
	return r.db.WithContext(c).
		Delete(&model.Friendship{}, "id = ?", friendshipID).Error
}

func (r *FriendshipRepository) Save(c context.Context, friendship *model.Friendship) error {
	return r.db.WithContext(c).Save(friendship).Error
}

func (r *FriendshipRepository) GetMessageByID(c context.Context, id string) (*model.Message, error) {
	var msg model.Message
	err := r.db.WithContext(c).First(&msg, "id = ?", id).Error
	return &msg, err
}


