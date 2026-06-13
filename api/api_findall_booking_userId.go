package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

// FindAllBookingByUserApi handles GET /FindallBookingByUser
// - Regular users: returns their own bookings (userId from JWT)
// - Admin/moderator: can pass ?UserId= to query any user's bookings
func FindAllBookingByUserApi(c *fiber.Ctx) error {

	jwtUserId, _ := c.Locals("userId").(string)
	role, _ := c.Locals("role").(string)

	targetUserId := c.Query("UserId")

	// Non-admin users can only see their own bookings
	if role != "admin" && role != "moderator" {
		targetUserId = jwtUserId
	} else if targetUserId == "" {
		targetUserId = jwtUserId
	}

	if targetUserId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "UserId is required")
	}

	fromDate := c.Query("FromDate")
	toDate := c.Query("ToDate")
	status := c.Query("Status")

	dateFilter := functions.BuildDateFilter("", fromDate, toDate)
	if status != "" {
		dateFilter["Status"] = status
	}

	bookings, err := dao.FindAllBookingsByUserId(targetUserId, dateFilter)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch bookings")
	}

	if bookings == nil {
		bookings = []dto.Booking{}
	}

	return utils.SendDataResponse(c, bookings)
}
