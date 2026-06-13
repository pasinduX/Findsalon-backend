package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func CountUnreadApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)
	if userId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "UserId is required")
	}

	count, err := dao.CountUnreadNotifications(userId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to count unread notifications")
	}

	return utils.SendDataResponse(c, fiber.Map{"UnreadCount": count})
}
