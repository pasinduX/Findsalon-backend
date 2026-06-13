package api

import (
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateSalonServiceApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)

	var service dto.SalonService
	if err := c.BodyParser(&service); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := validate.Struct(service); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if isOwner, _ := functions.IsSalonOwner(userId, service.SalonId); !isOwner {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Access denied: you are not the owner of this salon")
	}

	service.ServiceId = uuid.New().String()
	now := time.Now()
	service.CreatedAt = now
	service.UpdatedAt = now
	service.Deleted = false

	if err := dao.CreateSalonService(service); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create salon service")
	}
	return utils.SendSuccessResponse(c)
}
