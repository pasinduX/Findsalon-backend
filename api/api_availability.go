package api

import (
	"errors"
	"time"

	"findsalon-backend/dao"
	salonerr "findsalon-backend/errors"
	"findsalon-backend/functions"
	"findsalon-backend/repository"
	"findsalon-backend/service"
	"findsalon-backend/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"findsalon-backend/dto"
)

// bookingSvc is the shared service instance wired up at startup.
// In a larger app you'd use dependency injection; here a package-level var keeps it simple.
var bookingSvc *service.BookingService

func init() {
	bookingSvc = service.NewBookingService(
		&repository.BookingRepository{},
		&repository.ScheduleRepository{},
		&repository.BlockRepository{},
		&repository.ServiceRepository{},
	)
}

// GET /api/v1/booking/Availability?BarberId=&ServiceId=&Date=YYYY-MM-DD
func GetAvailabilityApi(c *fiber.Ctx) error {
	var req dto.AvailabilityRequest
	if err := c.QueryParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters")
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	slots, err := bookingSvc.GetAvailability(c.Context(), req)
	if err != nil {
		return availabilityError(c, err)
	}
	return utils.SendDataResponse(c, slots)
}

// POST /api/v1/booking/DirectBooking
func DirectBookingApi(c *fiber.Ctx) error {
	var req dto.DirectBookingRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	userId, _ := c.Locals("userId").(string)
	if userId == "" {
		userId = req.CustomerId
	}
	if userId == "" {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "Authenticated user is required")
	}

	user, err := dao.FindUserByUserId(userId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "Authenticated user is not synced")
	}

	req.CustomerId = user.UserId
	if req.CustomerName == "" {
		req.CustomerName = user.FullName
	}
	if req.CustomerPhone == "" {
		req.CustomerPhone = user.Phone
	}

	booking, err := bookingSvc.CreateDirectBooking(c.Context(), req)
	if err != nil {
		return availabilityError(c, err)
	}
	go functions.NotifyBooking(functions.BookingNotificationPayload{
		BookingId:     booking.BookingId,
		UserId:        booking.UserId,
		SalonId:       booking.SalonId,
		BarberId:      booking.BarberId,
		CustomerName:  booking.CustomerName,
		CustomerPhone: booking.CustomerPhone,
		Date:          booking.StartTime.Format("2006-01-02"),
		StartTime:     booking.StartTime.Format(time.Kitchen),
		EndTime:       booking.EndTime.Format(time.Kitchen),
		EventType:     dto.EventBookingCreated,
	})
	return utils.SendDataResponse(c, booking)
}

// availabilityError maps typed sentinel errors to HTTP status codes.
func availabilityError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, salonerr.ErrSlotTaken):
		return utils.SendErrorResponse(c, fiber.StatusConflict, "This slot is no longer available")
	case errors.Is(err, salonerr.ErrNotWorkingDay):
		return utils.SendErrorResponse(c, fiber.StatusUnprocessableEntity, "The barber does not work on this day")
	case errors.Is(err, salonerr.ErrOutsideWorkHours):
		return utils.SendErrorResponse(c, fiber.StatusUnprocessableEntity, "Requested time is outside working hours")
	case errors.Is(err, salonerr.ErrPastTime):
		return utils.SendErrorResponse(c, fiber.StatusUnprocessableEntity, "Cannot book a slot in the past")
	case errors.Is(err, salonerr.ErrNotFound):
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Resource not found")
	case errors.Is(err, salonerr.ErrInvalidInput):
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	case errors.Is(err, salonerr.ErrInvalidTimezone):
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	default:
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}
}
