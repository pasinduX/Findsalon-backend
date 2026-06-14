package api

import (
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func SendBookingNotificationApi(c *fiber.Ctx) error {
	var payload dto.BookingNotificationPayload
	if err := c.BodyParser(&payload); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if payload.BookingId == "" || payload.EventType == "" || payload.SalonId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "BookingId, EventType and SalonId are required")
	}

	go functions.NotifyBooking(functions.BookingNotificationPayload{
		BookingId:     payload.BookingId,
		UserId:        payload.UserId,
		SalonId:       payload.SalonId,
		BarberId:      payload.BarberId,
		CustomerName:  payload.CustomerName,
		CustomerEmail: payload.CustomerEmail,
		CustomerPhone: payload.CustomerPhone,
		Date:          payload.Date,
		StartTime:     payload.StartTime,
		EndTime:       payload.EndTime,
		EventType:     payload.EventType,
	})

	return utils.SendSuccessResponse(c)
}
