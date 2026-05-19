package handler

import (
	"chatapp/internal/dto"
	"chatapp/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func NewMessageHandler(svc *service.MessageService) *MessageHandler {
	return &MessageHandler{svc}
}

type MessageHandler struct {
	svc *service.MessageService
}

func (h *MessageHandler) GetConversations(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	convos, err := h.svc.GetConversations(c, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    convos,
	})
}

func (h *MessageHandler) GetMessages(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	conversationID := c.Params("conversationId")
	if conversationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "conversationId is required",
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

	offsetStr := c.Query("offset")
	offset := 0
	if offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil {
			if val >= 0 {
				offset = val
			}
		}
	}

	res, err := h.svc.GetMessages(c, userID, conversationID, limit, offset)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *MessageHandler) SendMessage(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	var req dto.SendMessageRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body format",
		})
	}

	if req.ReceiverID == "" || req.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "receiverId and text are required",
		})
	}

	msg, err := h.svc.SendMessage(c, userID, req.ReceiverID, req.Text)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    msg,
	})
}

func (h *MessageHandler) MarkMessageRead(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	messageID := c.Params("messageId")
	if messageID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "messageId is required",
		})
	}

	if err := h.svc.MarkAsRead(c, userID, messageID); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Message marked as read",
	})
}
