package api

import (
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateTemplateApi(c *fiber.Ctx) error {
	var payload dto.Template
	if err := c.BodyParser(&payload); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if payload.EventType != dto.EventBookingCreated && payload.EventType != dto.EventBookingCancelled && payload.EventType != dto.EventBookingCompleted && payload.EventType != dto.EventReviewReceived && payload.EventType != dto.EventCustom && payload.EventType != dto.EventBulk {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid EventType")
	}

	if !functions.UniqueCheck(dbConfig.TEMPLATES_COLLECTION, "EventType", payload.EventType) {
		return utils.SendErrorResponse(c, fiber.StatusConflict, "An active template for this event type already exists")
	}

	payload.TemplateId = uuid.New().String()
	payload.IsActive = true
	payload.Deleted = false
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()

	if err := dao.CreateTemplate(payload); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create template")
	}
	return utils.SendSuccessResponse(c)
}
