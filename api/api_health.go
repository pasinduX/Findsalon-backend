package api

import (
	"time"

	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func HealthApi(c *fiber.Ctx) error {
	return utils.SendDataResponse(c, fiber.Map{
		"service":   "findsalon-backend",
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
