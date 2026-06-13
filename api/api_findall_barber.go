package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func FindAllBarberApi(c *fiber.Ctx) error {
	salonId := c.Query("SalonId")
	userId := c.Query("UserId")

	if salonId == "" && userId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SalonId or UserId is required")
	}

	if userId != "" {
		// Return all barber records linked to this user (across all salons).
		barbers, err := dao.FindAllBarbersByUserId(userId)
		if err != nil {
			return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve barbers")
		}
		return utils.SendDataResponse(c, barbers)
	}

	barbers, err := dao.FindAllBarbersBySalonId(salonId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve barbers")
	}
	return utils.SendDataResponse(c, barbers)
}
