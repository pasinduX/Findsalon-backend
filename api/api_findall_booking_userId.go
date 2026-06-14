package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

// FindAllBookingByUserApi handles GET /FindallBookingByUser
// - When JWT auth is active:  regular users see only their own bookings (userId from JWT).
// - When JWT is not active:   falls back to ?UserId= query param.
// - Admin/moderator can pass ?UserId= to query any user's bookings.
func FindAllBookingByUserApi(c *fiber.Ctx) error {
	jwtUserId, _ := c.Locals("userId").(string)
	role, _ := c.Locals("role").(string)

	targetUserId := c.Query("UserId")

	if role == "admin" || role == "moderator" {
		// Admins may query any user; fall back to their own if none given.
		if targetUserId == "" {
			targetUserId = jwtUserId
		}
	} else if jwtUserId != "" {
		// JWT auth active — enforce own bookings only.
		targetUserId = jwtUserId
	}
	// else: JWT not configured, keep targetUserId from query param.

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

	// Enrich with salon and barber names, caching lookups by ID.
	salonCache := make(map[string]dto.Salon)
	barberCache := make(map[string]dto.Barber)

	enriched := make([]dto.BookingEnriched, 0, len(bookings))
	for _, b := range bookings {
		salon, ok := salonCache[b.SalonId]
		if !ok {
			salon, _ = dao.FindSalonBySalonId(b.SalonId)
			salonCache[b.SalonId] = salon
		}
		barber, ok := barberCache[b.BarberId]
		if !ok {
			barber, _ = dao.FindBarberByBarberId(b.BarberId)
			barberCache[b.BarberId] = barber
		}

		enriched = append(enriched, dto.BookingEnriched{
			BookingId:    b.BookingId,
			UserId:       b.UserId,
			SalonId:      b.SalonId,
			SalonName:    salon.Name,
			SalonAddress: salon.Address,
			SalonArea:    salon.Area,
			BarberId:     b.BarberId,
			BarberName:   barber.Name,
			SlotId:       b.SlotId,
			ServiceId:    b.ServiceId,
			Status:       b.Status,
			BookingType:  b.BookingType,
			CustomerName: b.CustomerName,
			Notes:        b.Notes,
			StartTime:    b.StartTime,
			EndTime:      b.EndTime,
			CreatedAt:    b.CreatedAt,
		})
	}

	return utils.SendDataResponse(c, fiber.Map{
		"Data": enriched,
	})
}
