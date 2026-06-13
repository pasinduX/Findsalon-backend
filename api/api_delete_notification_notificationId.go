package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func DeleteNotificationApi(c *fiber.Ctx) error {
	notificationId := c.Query("NotificationId")
	if notificationId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "NotificationId is required")
	}

	notification, err := dao.FindNotificationById(notificationId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Notification not found")
	}

	userId, _ := c.Locals("userId").(string)
	role, _ := c.Locals("role").(string)
	if notification.UserId != userId && role != "admin" {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Not authorized to delete this notification")
	}

	if err := dao.DeleteNotification(notificationId); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete notification")
	}

	return utils.SendSuccessResponse(c)
}
