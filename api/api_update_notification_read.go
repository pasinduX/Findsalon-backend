package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func MarkNotificationReadApi(c *fiber.Ctx) error {
	notificationId := c.Query("NotificationId")
	if notificationId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "NotificationId is required")
	}

	notification, err := dao.FindNotificationById(notificationId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Notification not found")
	}

	userId, _ := c.Locals("userId").(string)
	if notification.UserId != userId {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Not authorized to update this notification")
	}

	if err := dao.UpdateNotification(notificationId, bson.M{"IsRead": true}); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to mark notification as read")
	}

	return utils.SendSuccessResponse(c)
}
