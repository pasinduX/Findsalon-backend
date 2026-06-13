package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func FindAllTimeSlotApi(c *fiber.Ctx) error {
	salonId := c.Query("SalonId")
	barberId := c.Query("BarberId")
	date := c.Query("Date")
	fromDate := c.Query("FromDate")
	toDate := c.Query("ToDate")
	isBookedParam := c.Query("IsBooked")

	var slots interface{}
	var err error

	if date != "" && salonId != "" {
		slots, err = dao.FindAllTimeSlotsByDate(date, salonId)
	} else if barberId != "" {
		dateFilter := functions.BuildDateFilter(date, fromDate, toDate)
		if isBookedParam == "false" {
			dateFilter["IsBooked"] = false
		}
		slots, err = dao.FindAllTimeSlotsByBarberId(barberId, dateFilter)
	} else if salonId != "" {
		dateFilter := functions.BuildDateFilter(date, fromDate, toDate)
		if isBookedParam == "false" {
			dateFilter["IsBooked"] = false
		}
		slots, err = dao.FindAllTimeSlotsBySalonId(salonId, dateFilter)
	} else {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SalonId or BarberId is required")
	}

	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch time slots")
	}

	return utils.SendDataResponse(c, slots)
}
