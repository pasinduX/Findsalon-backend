package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func FindAllNotificationApi(c *fiber.Ctx) error {
	requestUserId := c.Query("UserId")
	currentUserId, _ := c.Locals("userId").(string)
	role, _ := c.Locals("role").(string)

	if requestUserId == "" {
		requestUserId = currentUserId
	}

	if requestUserId != currentUserId && role != "admin" && role != "moderator" {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Not authorized to view notifications for this user")
	}

	pagination := functions.GetPaginationParams(c)
	skip, limit := functions.BuildSkipLimit(pagination)
	notifications, total, err := dao.FindAllNotificationsByUserId(requestUserId, skip, limit)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to load notifications")
	}

	response := functions.BuildPaginatedResponse(notifications, total, pagination)
	return utils.SendDataResponse(c, response)
}
