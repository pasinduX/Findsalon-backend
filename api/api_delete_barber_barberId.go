package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func DeleteBarberApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)
	barberId := c.Query("BarberId")
	salonId := c.Query("SalonId")

	if barberId == "" || salonId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "BarberId and SalonId are required")
	}

	if isOwner, _ := functions.IsSalonOwner(userId, salonId); !isOwner {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Access denied: you are not the owner of this salon")
	}

	if err := dao.DeleteBarber(barberId); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete barber")
	}
	return utils.SendSuccessResponse(c)
}
