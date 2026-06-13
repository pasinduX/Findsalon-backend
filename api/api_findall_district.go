package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func FindAllDistrictApi(c *fiber.Ctx) error {
	districts, err := dao.FindAllDistricts()
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch districts")
	}
	return utils.SendDataResponse(c, districts)
}
