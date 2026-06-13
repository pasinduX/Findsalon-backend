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

type BulkSlotRequest struct {
	BarberId  string   `json:"BarberId" validate:"required"`
	SalonId   string   `json:"SalonId" validate:"required"`
	Dates     []string `json:"Dates" validate:"required"`
	StartTime string   `json:"StartTime" validate:"required"`
	EndTime   string   `json:"EndTime" validate:"required"`
}

func BulkCreateTimeSlotApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)

	var body BulkSlotRequest
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

	today := time.Now().Format("2006-01-02")
	collection := dbConfig.DATABASE.Collection(dbConfig.TIMESLOTS_COLLECTION)

	var validSlots []dto.TimeSlot
	skippedCount := 0

	for _, date := range body.Dates {
		if _, err := time.Parse("2006-01-02", date); err != nil {
			skippedCount++
			continue
		}
		if date < today {
			skippedCount++
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		count, err := collection.CountDocuments(ctx, bson.M{
			"BarberId":  body.BarberId,
			"Date":      date,
			"StartTime": body.StartTime,
			"Deleted":   false,
		})
		cancel()

		if err != nil || count > 0 {
			skippedCount++
			continue
		}

		now := time.Now()
		validSlots = append(validSlots, dto.TimeSlot{
			SlotId:    uuid.New().String(),
			BarberId:  body.BarberId,
			SalonId:   body.SalonId,
			Date:      date,
			StartTime: body.StartTime,
			EndTime:   body.EndTime,
			IsBooked:  false,
			CreatedAt: now,
			UpdatedAt: now,
			Deleted:   false,
		})
	}

	createdCount := len(validSlots)
	if createdCount > 0 {
		if err := dao.BulkCreateTimeSlots(validSlots); err != nil {
			return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to bulk create time slots")
		}
	}

	return utils.SendDataResponse(c, fiber.Map{
		"Created": createdCount,
		"Skipped": skippedCount,
	})
}
