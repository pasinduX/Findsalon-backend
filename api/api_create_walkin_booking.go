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

type WalkInRequest struct {
	SalonId      string `json:"SalonId" validate:"required"`
	BarberId     string `json:"BarberId" validate:"required"`
	SlotId       string `json:"SlotId" validate:"required"`
	CustomerName string `json:"CustomerName" validate:"required"`
	Notes        string `json:"Notes"`
}

func CreateWalkInBookingApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)

	var body WalkInRequest
	if err := c.BodyParser(&body); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	isAllowed, _ := functions.IsSalonOwnerOrBarber(userId, body.SalonId)
	if !isAllowed {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "You are not authorized to create walk-in bookings for this salon")
	}

	available, err := functions.IsSlotAvailable(body.SlotId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to check slot availability")
	}
	if !available {
		return utils.SendErrorResponse(c, fiber.StatusConflict, "Time slot is not available")
	}

	now := time.Now()
	booking := dto.Booking{
		BookingId:    uuid.New().String(),
		UserId:       "",
		SalonId:      body.SalonId,
		BarberId:     body.BarberId,
		SlotId:       body.SlotId,
		Status:       dto.BookingStatusConfirmed,
		BookingType:  dto.BookingTypeWalkIn,
		CustomerName: body.CustomerName,
		Notes:        body.Notes,
		CreatedAt:    now,
		UpdatedAt:    now,
		Deleted:      false,
	}

	if err := dao.CreateBooking(booking); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create walk-in booking")
	}

	if err := dao.MarkTimeSlotBooked(body.SlotId); err != nil {
		log.Printf("CreateWalkInBookingApi: failed to mark slot booked: %v", err)
	}

	go functions.NotifyBooking(functions.BookingNotificationPayload{
		BookingId:    booking.BookingId,
		UserId:       booking.UserId,
		SalonId:      booking.SalonId,
		BarberId:     booking.BarberId,
		CustomerName: booking.CustomerName,
		EventType:    "booking_created",
	})

	return utils.SendSuccessResponse(c)
}
