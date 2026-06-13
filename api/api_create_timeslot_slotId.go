package api

import (
	"context"
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateTimeSlotApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)

	var body dto.TimeSlot
	if err := c.BodyParser(&body); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	isAllowed, err := functions.IsSalonOwnerOrBarber(userId, body.SalonId)
	if err != nil || !isAllowed {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "You are not authorized to create slots for this salon")
	}

	_, err = time.Parse("2006-01-02", body.Date)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid date format, expected YYYY-MM-DD")
	}
	if body.Date < time.Now().Format("2006-01-02") {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Date cannot be in the past")
	}

	_, err = time.Parse("15:04", body.StartTime)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid StartTime format, expected HH:MM")
	}
	_, err = time.Parse("15:04", body.EndTime)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid EndTime format, expected HH:MM")
	}
	if body.StartTime >= body.EndTime {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "StartTime must be before EndTime")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.TIMESLOTS_COLLECTION)
	count, err := collection.CountDocuments(ctx, bson.M{
		"BarberId":  body.BarberId,
		"Date":      body.Date,
		"StartTime": body.StartTime,
		"Deleted":   false,
	})
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to check duplicate slot")
	}
	if count > 0 {
		return utils.SendErrorResponse(c, fiber.StatusConflict, "Time slot already exists for this barber at this time")
	}

	now := time.Now()
	body.SlotId = uuid.New().String()
	body.IsBooked = false
	body.Deleted = false
	body.CreatedAt = now
	body.UpdatedAt = now

	if err := dao.CreateTimeSlot(body); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create time slot")
	}

	return utils.SendSuccessResponse(c)
}
