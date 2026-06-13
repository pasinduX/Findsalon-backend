package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func FindAllBookingApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)
	role, _ := c.Locals("role").(string)

	salonId := c.Query("SalonId")
	barberId := c.Query("BarberId")
	queryUserId := c.Query("UserId")
	date := c.Query("Date")
	fromDate := c.Query("FromDate")
	toDate := c.Query("ToDate")
	status := c.Query("Status")

	dateFilter := functions.BuildDateFilter(date, fromDate, toDate)
	if status != "" {
		dateFilter["Status"] = status
	}

	var bookings []dto.Booking
	var err error

	switch role {
	case "admin", "moderator":
		if date != "" && salonId != "" {
			bookings, err = dao.FindAllBookingsByDate(date, salonId)
		} else if salonId != "" {
			bookings, err = dao.FindAllBookingsBySalonId(salonId, dateFilter)
		} else if barberId != "" {
			bookings, err = dao.FindAllBookingsByBarberId(barberId, dateFilter)
		} else if queryUserId != "" {
			bookings, err = dao.FindAllBookingsByUserId(queryUserId, dateFilter)
		} else {
			bookings, err = dao.FindAllBookingsBySalonId("", dateFilter)
		}

	case "salon_owner":
		if salonId == "" {
			return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SalonId is required for salon owners")
		}
		isOwner, _ := functions.IsSalonOwner(userId, salonId)
		if !isOwner {
			return utils.SendErrorResponse(c, fiber.StatusForbidden, "You are not the owner of this salon")
		}
		if date != "" {
			bookings, err = dao.FindAllBookingsByDate(date, salonId)
		} else {
			bookings, err = dao.FindAllBookingsBySalonId(salonId, dateFilter)
		}

	case "barber":
		if salonId == "" {
			return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SalonId is required for barbers")
		}
		isBarber, _ := functions.IsBarberOfSalon(userId, salonId)
		if !isBarber {
			return utils.SendErrorResponse(c, fiber.StatusForbidden, "You are not a barber of this salon")
		}
		if barberId == "" {
			barberId = userId
		}
		bookings, err = dao.FindAllBookingsByBarberId(barberId, dateFilter)

	default:
		forcedFilter := bson.M{}
		for k, v := range dateFilter {
			forcedFilter[k] = v
		}
		bookings, err = dao.FindAllBookingsByUserId(userId, forcedFilter)
	}

	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch bookings")
	}

	if bookings == nil {
		bookings = []dto.Booking{}
	}

	return utils.SendDataResponse(c, bookings)
}
