package internal

import (
	"chatapp/internal/handler"
	"chatapp/internal/middleware"
	"chatapp/internal/repository"
	"chatapp/internal/service"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// 1. Initialize repositories
	userRepo := repository.NewUserRepository(db)
	friendshipRepo := repository.NewFriendshipRepository(db)
	convoRepo := repository.NewConversationRepository(db)
	messageRepo := repository.NewMessageRepository(db)

	// 2. Initialize services
	authService := service.NewAuthService(userRepo)
	friendService := service.NewFriendService(friendshipRepo, userRepo, convoRepo)
	messageService := service.NewMessageService(messageRepo, convoRepo, userRepo, friendshipRepo)

	// 3. Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userRepo)
	friendHandler := handler.NewFriendHandler(friendService)
	messageHandler := handler.NewMessageHandler(messageService)

	// 4. Group all routes under /api/v1 prefix
	api := app.Group("/api/v1")

	// Authentication routes
	api.Post("/auth/register", authHandler.Register)
	api.Post("/auth/login", authHandler.Login)
	api.Delete("/auth/logout", middleware.AuthProtected(), authHandler.Logout)
	api.Get("/auth/me", middleware.AuthProtected(), authHandler.GetMe)

	// User routes
	api.Get("/users/search", middleware.AuthProtected(), userHandler.SearchUsers)
	api.Get("/users/:userId", middleware.AuthProtected(), userHandler.GetUserByID)

	// Friend routes
	api.Get("/friends", middleware.AuthProtected(), friendHandler.GetFriends)
	api.Get("/friends/requests", middleware.AuthProtected(), friendHandler.GetFriendRequests)
	api.Post("/friends/request", middleware.AuthProtected(), friendHandler.SendFriendRequest)
	api.Post("/friends/accept", middleware.AuthProtected(), friendHandler.AcceptFriendRequest)
	api.Post("/friends/decline", middleware.AuthProtected(), friendHandler.DeclineFriendRequest)
	api.Delete("/friends/:friendId", middleware.AuthProtected(), friendHandler.RemoveFriend)

	// Message routes
	api.Get("/messages/conversations", middleware.AuthProtected(), messageHandler.GetConversations)
	api.Get("/messages/:conversationId", middleware.AuthProtected(), messageHandler.GetMessages)
	api.Post("/messages", middleware.AuthProtected(), messageHandler.SendMessage)
	api.Patch("/messages/:messageId/read", middleware.AuthProtected(), messageHandler.MarkMessageRead)
}
