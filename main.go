package main

import (
	"log"
	"os"
	"path/filepath"
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
	loadEnvironment()

	integrations.SetEnvironmentVariables()
	dbConfig.ConnectToMongoDB()

	if err := functions.SeedDistricts(); err != nil {
		log.Printf("Warning: failed to seed districts: %v", err)
	}
	if err := functions.SeedSpecialties(); err != nil {
		log.Printf("Warning: failed to seed specialties: %v", err)
	}
	functions.SeedDefaultTemplates()

	app := fiber.New(fiber.Config{
		AppName:        "FindSalon v1.0",
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		ReadBufferSize: 16384, // 16KB — Auth0 session cookies on localhost can exceed the 4KB default
	})

	app.Use(logger.New())
	app.Use(recover.New())

	apiHandlers.Router(app)

	port := integrations.ServerPort
	if port == "" {
		port = "8080"
	}
	log.Printf("FindSalon backend starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func loadEnvironment() {
	if envPath := findEnvFile(); envPath != "" {
		if err := godotenv.Load(envPath); err != nil {
			log.Printf("Warning: failed to load %s: %v", envPath, err)
		}
		return
	}

	wd, _ := os.Getwd()
	log.Printf("Warning: .env file not found from %s, reading from environment", wd)
}

func findEnvFile() string {
	seen := map[string]bool{}
	for _, start := range envSearchRoots() {
		for dir := start; dir != ""; dir = filepath.Dir(dir) {
			if seen[dir] {
				break
			}
			seen[dir] = true

			envPath := filepath.Join(dir, ".env")
			if info, err := os.Stat(envPath); err == nil && !info.IsDir() {
				return envPath
			}

			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
		}
	}
	return ""
}

func envSearchRoots() []string {
	roots := make([]string, 0, 2)
	if wd, err := os.Getwd(); err == nil {
		roots = append(roots, wd)
	}
	if exe, err := os.Executable(); err == nil {
		roots = append(roots, filepath.Dir(exe))
	}
	return roots
}
