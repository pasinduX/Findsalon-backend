package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func MarkAllReadApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)
	if userId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "UserId is required")
	}

	if err := dao.MarkAllNotificationsRead(userId); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to mark all notifications as read")
	}
	return utils.SendSuccessResponse(c)
}
