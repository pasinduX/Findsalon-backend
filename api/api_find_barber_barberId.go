package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func FindBarberApi(c *fiber.Ctx) error {
	barberId := c.Query("BarberId")
	if barberId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "BarberId is required")
	}

	barber, err := dao.FindBarberByBarberId(barberId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Barber not found")
	}
	return utils.SendDataResponse(c, barber)
}
