package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func FindBookingApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)

	bookingId := c.Query("BookingId")
	if bookingId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "BookingId is required")
	}

	booking, err := dao.FindBookingByBookingId(bookingId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Booking not found")
	}

	isOwner := booking.UserId == userId
	isSalonOwner, _ := functions.IsSalonOwner(userId, booking.SalonId)
	isBarber, _ := functions.IsBarberOfSlot(userId, booking.SlotId)

	if !isOwner && !isSalonOwner && !isBarber {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "You are not authorized to view this booking")
	}

	return utils.SendDataResponse(c, booking)
}
