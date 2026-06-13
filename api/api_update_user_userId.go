package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateUserApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)

	var updateMap bson.M
	if err := c.BodyParser(&updateMap); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	delete(updateMap, "PasswordHash")
	delete(updateMap, "RefreshToken")
	delete(updateMap, "UserId")
	delete(updateMap, "CreatedAt")
	delete(updateMap, "Deleted")

	if err := dao.UpdateUser(userId, updateMap); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to update user")
	}

	return utils.SendSuccessResponse(c)
}
