package api

import (
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func DownloadBookingApi(c *fiber.Ctx) error {
	role, _ := c.Locals("role").(string)

	if role != "admin" && role != "moderator" {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Only admin or moderator can download bookings")
	}

	salonId := c.Query("SalonId")
	fromDate := c.Query("FromDate")
	toDate := c.Query("ToDate")

	dateFilter := functions.BuildDateFilter("", fromDate, toDate)

	bookings, err := dao.FindAllBookingsBySalonId(salonId, dateFilter)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch bookings")
	}

	filename := "bookings_export_" + time.Now().Format("20060102")
	return utils.SendJSONFileDownload(c, bookings, filename)
}
