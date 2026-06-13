package main

import (
	"log"
	"os"
	"time"

	"findsalon-backend/apiHandlers"
	"findsalon-backend/dbConfig"
	"findsalon-backend/functions"
	"findsalon-backend/integrations"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found, reading from environment")
	}

	integrations.SetEnvironmentVariables()
	dbConfig.ConnectToMongoDB()

	if err := functions.SeedDistricts(); err != nil {
		log.Printf("Warning: failed to seed districts: %v", err)
	}
	if err := functions.SeedSpecialties(); err != nil {
		log.Printf("Warning: failed to seed specialties: %v", err)
	}
	functions.SeedDefaultTemplates()

	storagePath := os.Getenv("STORAGE_BASE_PATH")
	if storagePath == "" {
		storagePath = "./uploads"
	}
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName:        "FindSalon v1.0",
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		ReadBufferSize: 16384, // 16KB — Auth0 session cookies on localhost can exceed the 4KB default
	})

	app.Use(logger.New())
	app.Use(recover.New())

	// Serve uploaded files (images) at /uploads/<path>
	app.Static("/uploads", storagePath)

	apiHandlers.Router(app)

	port := integrations.ServerPort
	if port == "" {
		port = "8080"
	}
	log.Printf("FindSalon backend starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
