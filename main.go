package main

import (
	"chatapp/database"
	"chatapp/internal"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
)

func GetFiberServer() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			var fiberErr *fiber.Error
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			if errors.As(err, &fiberErr) {
				code = fiberErr.Code
				message = fiberErr.Message
			} else {
				fmt.Println("Unhandled error:", err)
				message = err.Error()
			}

			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error":   message,
			})
		},
	})

	return app
}

func main() {
	seedFlag := flag.Bool("seed", false, "Seed the database with mock data")
	flag.Parse()

	db := database.GetDBConnection()

	if *seedFlag {
		log.Println("Seeding database...")
		if err := database.Seed(db); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
		log.Println("Database seeded successfully!")
		return
	}

	app := GetFiberServer()


	internal.SetupRoutes(app, db)
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "3000"
		}
		if err := app.Listen(":" + port); err != nil {
			log.Printf("Failed to start server : %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Failed to shutdown app : %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Gagal mendapatkan instance DB: %v", err)
	} else {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Gagal menutup database: %v", err)
		}
	}

	log.Println("Sistem shutdown gracefully")
}
