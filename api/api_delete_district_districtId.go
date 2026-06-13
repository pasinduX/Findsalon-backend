package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func DeleteDistrictApi(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "District id is required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid district id")
	}

	if err := dao.DeleteDistrict(id); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete district")
	}
	return utils.SendSuccessResponse(c)
}
