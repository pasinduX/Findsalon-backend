package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func FindAllQuoteApi(c *fiber.Ctx) error {
	salonId := c.Query("SalonId")
	if salonId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SalonId is required")
	}

	quotes, err := dao.FindAllQuotesBySalonId(salonId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve quotes")
	}
	return utils.SendDataResponse(c, quotes)
}
