package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func DeleteWorkingHoursApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)
	hoursId := c.Query("HoursId")
	salonId := c.Query("SalonId")

	if hoursId == "" || salonId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "HoursId and SalonId are required")
	}

	if isOwner, _ := functions.IsSalonOwner(userId, salonId); !isOwner {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Access denied: you are not the owner of this salon")
	}

	if err := dao.DeleteWorkingHours(hoursId); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete working hours")
	}
	return utils.SendSuccessResponse(c)
}
