package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func CreateDistrictApi(c *fiber.Ctx) error {
	var district dto.DistrictDTO
	if err := c.BodyParser(&district); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if district.ID == 0 {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "District ID is required")
	}
	if district.Name == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "District name is required")
	}

	if district.Cities == nil {
		district.Cities = []dto.CityDTO{}
	}

	if err := dao.CreateDistrict(district); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create district")
	}
	return utils.SendSuccessResponse(c)
}
