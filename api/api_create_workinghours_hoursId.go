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

func CreateWorkingHoursApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)

	var hours dto.WorkingHours
	if err := c.BodyParser(&hours); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := validate.Struct(hours); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if isOwner, _ := functions.IsSalonOwner(userId, hours.SalonId); !isOwner {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Access denied: you are not the owner of this salon")
	}

	hours.HoursId = uuid.New().String()
	now := time.Now()
	hours.CreatedAt = now
	hours.UpdatedAt = now
	hours.Deleted = false

	if err := dao.CreateWorkingHours(hours); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create working hours")
	}
	return utils.SendSuccessResponse(c)
}
