package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
)

func UpdateDistrictApi(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "District id is required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid district id")
	}

	var update bson.M
	if err := c.BodyParser(&update); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	delete(update, "id")

	if len(update) == 0 {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "No fields to update")
	}

	if err := dao.UpdateDistrict(id, update); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to update district")
	}
	return utils.SendSuccessResponse(c)
}
