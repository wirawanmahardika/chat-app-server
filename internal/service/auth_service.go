package service

import (
	"chatapp/internal/dto"
	"chatapp/internal/repository"
	"chatapp/pkg/utils"
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo}
}

type AuthService struct {
	userRepository *repository.UserRepository
}

func (s *AuthService) Register(c context.Context, name, email, password string) (*dto.AuthResponse, error) {
	count, err := s.userRepository.CountByEmail(c, email)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, fiber.NewError(fiber.StatusConflict, "Email already exists")
	}

	user, err := s.userRepository.Create(c, name, email, password)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to generate authentication token")
	}

	// Update status
	_ = s.userRepository.UpdateOnlineStatus(c, user.ID, true)

	return &dto.AuthResponse{
		Success: true,
		Token:   token,
		User: dto.UserProfileResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Avatar:    user.Avatar,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *AuthService) Login(c context.Context, email, password string) (*dto.AuthResponse, error) {
	user, err := s.userRepository.GetByEmail(c, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
		}
		return nil, err
	}

	if !user.ValidatePassword(password) {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to generate authentication token")
	}

	// Update status
	_ = s.userRepository.UpdateOnlineStatus(c, user.ID, true)

	return &dto.AuthResponse{
		Success: true,
		Token:   token,
		User: dto.UserProfileResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Avatar:    user.Avatar,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *AuthService) Logout(c context.Context, userID string) error {
	return s.userRepository.UpdateOnlineStatus(c, userID, false)
}

func (s *AuthService) GetProfile(c context.Context, userID string) (*dto.UserProfileResponse, error) {
	user, err := s.userRepository.GetByID(c, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return nil, err
	}

	return &dto.UserProfileResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}
