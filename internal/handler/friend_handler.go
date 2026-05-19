package handler

import (
	"chatapp/internal/dto"
	"chatapp/internal/service"

	"github.com/gofiber/fiber/v3"
)

func NewFriendHandler(svc *service.FriendService) *FriendHandler {
	return &FriendHandler{svc}
}

type FriendHandler struct {
	svc *service.FriendService
}

func (h *FriendHandler) GetFriends(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	friends, err := h.svc.GetFriends(c, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    friends,
	})
}

func (h *FriendHandler) GetFriendRequests(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	requests, err := h.svc.GetPendingRequests(c, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    requests,
	})
}

func (h *FriendHandler) SendFriendRequest(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	var req dto.SendFriendRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body format",
		})
	}

	if req.FriendID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "friendId is required",
		})
	}

	if err := h.svc.SendRequest(c, userID, req.FriendID); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Friend request sent",
	})
}

func (h *FriendHandler) AcceptFriendRequest(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	var req dto.RespondFriendRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body format",
		})
	}

	if req.RequestID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "requestId is required",
		})
	}

	friend, err := h.svc.AcceptRequest(c, userID, req.RequestID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Friend request accepted",
		"data":    friend,
	})
}

func (h *FriendHandler) DeclineFriendRequest(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	var req dto.RespondFriendRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body format",
		})
	}

	if req.RequestID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "requestId is required",
		})
	}

	if err := h.svc.DeclineRequest(c, userID, req.RequestID); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Friend request declined",
	})
}

func (h *FriendHandler) RemoveFriend(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	friendID := c.Params("friendId")
	if friendID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "friendId parameter is required",
		})
	}

	if err := h.svc.RemoveFriend(c, userID, friendID); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Friend removed",
	})
}
