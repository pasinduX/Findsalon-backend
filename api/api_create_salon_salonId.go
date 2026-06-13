package api

import (
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var validate = validator.New()

func CreateSalonApi(c *fiber.Ctx) error {
	var salon dto.Salon
	if err := c.BodyParser(&salon); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Auth middleware sets userId when security is enabled; fall back to body's OwnerId otherwise
	if localUserId, ok := c.Locals("userId").(string); ok && localUserId != "" {
		salon.OwnerId = localUserId
	}
	if salon.OwnerId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "OwnerId is required")
	}
	if err := validate.Struct(salon); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	salon.SalonId = uuid.New().String()
	now := time.Now()
	salon.CreatedAt = now
	salon.UpdatedAt = now
	salon.IsActive = true
	salon.Deleted = false

	// Mirror lat/lng into a GeoJSON Point for 2dsphere geo queries.
	if salon.Location.Latitude != 0 || salon.Location.Longitude != 0 {
		salon.GeoLocation = dto.NewGeoPoint(salon.Location.Latitude, salon.Location.Longitude)
	}

	if err := dao.CreateSalon(salon); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create salon")
	}
	return utils.SendSuccessResponse(c)
}
