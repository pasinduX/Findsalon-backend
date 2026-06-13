package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func DeleteBookingApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)

	bookingId := c.Query("BookingId")
	if bookingId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "BookingId is required")
	}

	booking, err := dao.FindBookingByBookingId(bookingId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Booking not found")
	}

	if booking.Status != dto.BookingStatusCancelled {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Cancel the booking before deleting")
	}

	isOwner := booking.UserId == userId
	isSalonOwner, _ := functions.IsSalonOwner(userId, booking.SalonId)

	if !isOwner && !isSalonOwner {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "You are not authorized to delete this booking")
	}

	if err := dao.DeleteBooking(bookingId); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete booking")
	}

	return utils.SendSuccessResponse(c)
}
