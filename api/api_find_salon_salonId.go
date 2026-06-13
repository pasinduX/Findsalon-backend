package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func FindSalonApi(c *fiber.Ctx) error {
	salonId := c.Query("SalonId")
	if salonId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SalonId is required")
	}

	salon, err := dao.FindSalonBySalonId(salonId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Salon not found")
	}
	return utils.SendDataResponse(c, salon)
}
