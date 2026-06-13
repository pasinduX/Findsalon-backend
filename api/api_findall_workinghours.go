package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func FindAllWorkingHoursApi(c *fiber.Ctx) error {
	salonId := c.Query("SalonId")
	if salonId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SalonId is required")
	}

	hours, err := dao.FindAllWorkingHoursBySalonId(salonId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve working hours")
	}
	return utils.SendDataResponse(c, hours)
}
