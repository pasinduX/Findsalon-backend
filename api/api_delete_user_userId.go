package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func DeleteUserApi(c *fiber.Ctx) error {
	userId := c.Query("UserId")
	if userId == "" {
		userId, _ = c.Locals("userId").(string)
	}

	if err := dao.DeleteUser(userId); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete user")
	}

	return utils.SendSuccessResponse(c)
}
