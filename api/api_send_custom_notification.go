package api

import (
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SendCustomNotificationApi(c *fiber.Ctx) error {
	var payload dto.CustomNotificationPayload
	if err := c.BodyParser(&payload); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	v := validator.New()
	if err := v.Struct(payload); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	eventType := payload.EventType
	if eventType == "" {
		eventType = dto.EventCustom
	}

	notification := dto.Notification{
		NotificationId: uuid.New().String(),
		UserId:         payload.UserId,
		Title:          payload.Title,
		Body:           payload.Body,
		Type:           dto.NotificationTypeSystem,
		EventType:      eventType,
		RefId:          payload.RefId,
		IsRead:         false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Deleted:        false,
	}

	if err := dao.CreateNotification(notification); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create custom notification")
	}

	return utils.SendSuccessResponse(c)
}
