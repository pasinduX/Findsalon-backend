package api

import (
	"log"
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateBookingApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)

	var body dto.Booking
	if err := c.BodyParser(&body); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if userId != "" {
		body.UserId = userId
	}
	body.BookingType = dto.BookingTypeOnline

	available, err := functions.IsSlotAvailable(body.SlotId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to check slot availability")
	}
	if !available {
		return utils.SendErrorResponse(c, fiber.StatusConflict, "Time slot is not available")
	}

	valid, err := functions.ValidateSlotOwnership(body.SlotId, body.BarberId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to validate slot ownership")
	}
	if !valid {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SlotId does not belong to the specified BarberId")
	}

	now := time.Now()
	body.BookingId = uuid.New().String()
	body.Status = dto.BookingStatusConfirmed
	body.CreatedAt = now
	body.UpdatedAt = now
	body.Deleted = false

	if err := dao.CreateBooking(body); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create booking")
	}

	if err := dao.MarkTimeSlotBooked(body.SlotId); err != nil {
		log.Printf("CreateBookingApi: failed to mark slot booked: %v", err)
	}

	go functions.NotifyBooking(functions.BookingNotificationPayload{
		BookingId:     body.BookingId,
		UserId:        body.UserId,
		SalonId:       body.SalonId,
		BarberId:      body.BarberId,
		CustomerName:  body.CustomerName,
		CustomerPhone: body.CustomerPhone,
		Date:          body.StartTime.Format("2006-01-02"),
		StartTime:     body.StartTime.Format(time.Kitchen),
		EndTime:       body.EndTime.Format(time.Kitchen),
		EventType:     dto.EventBookingCreated,
	})

	return utils.SendSuccessResponse(c)
}
