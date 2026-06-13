package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func FindAllSalonServiceApi(c *fiber.Ctx) error {
	salonId := c.Query("SalonId")
	if salonId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SalonId is required")
	}

	services, err := dao.FindAllServicesBySalonId(salonId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve salon services")
	}
	return utils.SendDataResponse(c, services)
}
