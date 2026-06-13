package api

import (
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SendBulkNotificationApi(c *fiber.Ctx) error {
	var payload dto.BulkNotificationPayload
	if err := c.BodyParser(&payload); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	v := validator.New()
	if err := v.Struct(payload); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if len(payload.UserIds) == 0 {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "UserIds must not be empty")
	}
	if len(payload.UserIds) > 500 {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Bulk limit exceeded: max 500 recipients")
	}

	sent := 0
	for _, userId := range payload.UserIds {
		notification := dto.Notification{
			NotificationId: uuid.New().String(),
			UserId:         userId,
			Title:          payload.Title,
			Body:           payload.Body,
			Type:           dto.NotificationTypeSystem,
			EventType:      dto.EventBulk,
			IsRead:         false,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			Deleted:        false,
		}
		if err := dao.CreateNotification(notification); err != nil {
			return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create bulk notification")
		}
		sent++
		go func(uid string) {
			if email, _, err := functions.FetchUserEmailAndName(uid); err == nil && email != "" {
				functions.SendEmailAsync(email, payload.Title, "<p>"+payload.Body+"</p>")
			}
		}(userId)
	}

	return utils.SendDataResponse(c, fiber.Map{"Sent": sent})
}
