package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func FindUserRoleApi(c *fiber.Ctx) error {

	userId := c.Query("UserId")
	currentUserId, _ := c.Locals("userId").(string)
	role, _ := c.Locals("role").(string)
	if userId == "" {
		userId = currentUserId
	}
	if userId != currentUserId && role != "admin" && role != "moderator" {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Cannot view other users roles")
	}

	roles, err := dao.FindUserRolesByUserId(userId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to fetch user roles")
	}
	return utils.SendDataResponse(c, roles)
}
