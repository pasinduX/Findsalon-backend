package api

import (
	"log"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateBookingApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)
	role, _ := c.Locals("role").(string)

	bookingId := c.Query("BookingId")
	if bookingId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "BookingId is required")
	}

	booking, err := dao.FindBookingByBookingId(bookingId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Booking not found")
	}

	if booking.Status == dto.BookingStatusCancelled {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Cannot update a cancelled booking")
	}
	if booking.Status == dto.BookingStatusCompleted {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Cannot update a completed booking")
	}

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	updateFields := bson.M{}
	if notes, ok := body["Notes"]; ok {
		updateFields["Notes"] = notes
	}

	newStatus, hasStatus := body["Status"]
	if hasStatus {
		newStatusStr, ok := newStatus.(string)
		if !ok {
			return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid Status value")
		}

		isOwner := booking.UserId == userId
		isSalonOwner, _ := functions.IsSalonOwner(userId, booking.SalonId)
		isBarber, _ := functions.IsBarberOfSlot(userId, booking.SlotId)
		isAdminOrMod := role == "admin" || role == "moderator"

		switch newStatusStr {
		case dto.BookingStatusCancelled:
			if !isOwner && !isSalonOwner && !isBarber && !isAdminOrMod {
				return utils.SendErrorResponse(c, fiber.StatusForbidden, "Not authorized to cancel this booking")
			}
			updateFields["Status"] = dto.BookingStatusCancelled

			if err := dao.MarkTimeSlotAvailable(booking.SlotId); err != nil {
				log.Printf("UpdateBookingApi: failed to free slot: %v", err)
			}

			go functions.NotifyBooking(functions.BookingNotificationPayload{
				BookingId:    booking.BookingId,
				UserId:       booking.UserId,
				SalonId:      booking.SalonId,
				BarberId:     booking.BarberId,
				CustomerName: booking.CustomerName,
				Date:         booking.StartTime.Format("2006-01-02"),
				StartTime:    booking.StartTime.Format("3:04 PM"),
				EndTime:      booking.EndTime.Format("3:04 PM"),
				EventType:    dto.EventBookingCancelled,
			})

		case dto.BookingStatusCompleted:
			if !isSalonOwner && !isBarber && !isAdminOrMod {
				return utils.SendErrorResponse(c, fiber.StatusForbidden, "Only salon owner or barber can complete a booking")
			}
			updateFields["Status"] = dto.BookingStatusCompleted

			go functions.NotifyBooking(functions.BookingNotificationPayload{
				BookingId:    booking.BookingId,
				UserId:       booking.UserId,
				SalonId:      booking.SalonId,
				BarberId:     booking.BarberId,
				CustomerName: booking.CustomerName,
				Date:         booking.StartTime.Format("2006-01-02"),
				StartTime:    booking.StartTime.Format("3:04 PM"),
				EndTime:      booking.EndTime.Format("3:04 PM"),
				EventType:    dto.EventBookingCompleted,
			})

		default:
			return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid status transition")
		}
	}

	if len(updateFields) == 0 {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "No valid fields to update")
	}

	if err := dao.UpdateBooking(bookingId, bson.M{"$set": updateFields}); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to update booking")
	}

	return utils.SendSuccessResponse(c)
}
