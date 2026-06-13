package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func FindDistrictApi(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "District id is required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid district id")
	}

	district, err := dao.FindDistrictById(id)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "District not found")
	}
	return utils.SendDataResponse(c, district)
}
