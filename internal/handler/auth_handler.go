package handler

import (
	"chatapp/internal/dto"
	"chatapp/internal/service"

	"github.com/gofiber/fiber/v3"
)

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc}
}

type AuthHandler struct {
	svc *service.AuthService
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	req := dto.RegisterRequest{}

	println("ok")

	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body format",
		})
	}

	// Simple validations
	if req.Name == "" || req.Email == "" || len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Name and Email are required; Password must be at least 6 characters",
		})
	}

	res, err := h.svc.Register(c, req.Name, req.Email, req.Password)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	req := dto.LoginRequest{}

	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body format",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Email and Password are required",
		})
	}

	res, err := h.svc.Login(c, req.Email, req.Password)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	if err := h.svc.Logout(c, userID); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
	})
}

func (h *AuthHandler) GetMe(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	profile, err := h.svc.GetProfile(c, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    profile,
	})
}
