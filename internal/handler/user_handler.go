package handler

import (
	"chatapp/internal/repository"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo}
}

type UserHandler struct {
	userRepo *repository.UserRepository
}

func (h *UserHandler) SearchUsers(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	q := c.Query("q")
	if len(q) < 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Search query 'q' must be at least 2 characters",
		})
	}

	limitStr := c.Query("limit")
	limit := 50
	if limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil {
			if val > 0 && val <= 100 {
				limit = val
			}
		}
	}

	users, err := h.userRepo.Search(c, q, limit, userID)
	if err != nil {
		return err
	}

	// Format data structure matching swagger.yaml
	type UserResponseItem struct {
		ID        string  `json:"id"`
		Name      string  `json:"name"`
		Email     string  `json:"email"`
		Avatar    *string `json:"avatar"`
		Online    bool    `json:"online"`
		LastSeen  string  `json:"lastSeen"`
		CreatedAt string  `json:"createdAt"`
		UpdatedAt string  `json:"updatedAt"`
	}

	var res []UserResponseItem
	for _, u := range users {
		res = append(res, UserResponseItem{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Avatar:    u.Avatar,
			Online:    u.Online,
			LastSeen:  u.LastSeen.Format("2006-01-02T15:04:05Z"),
			CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	if res == nil {
		res = []UserResponseItem{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
	})
}

func (h *UserHandler) GetUserByID(c fiber.Ctx) error {
	targetUserID := c.Params("userId")
	if targetUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "User ID is required",
		})
	}

	user, err := h.userRepo.GetByID(c, targetUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "User not found",
			})
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":        user.ID,
			"name":      user.Name,
			"email":     user.Email,
			"avatar":    user.Avatar,
			"online":    user.Online,
			"lastSeen":  user.LastSeen.Format("2006-01-02T15:04:05Z"),
			"createdAt": user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			"updatedAt": user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
	})
}
