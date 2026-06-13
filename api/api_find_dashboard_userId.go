package api

import (
	"sync"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func BuildUserDashboard(userId string, authHeader string) (dto.UserDashboard, error) {
	var profile dto.Profile
	var bookings []dto.BookingSummary
	var profileErr error
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		profile, profileErr = dao.FindProfileByUserId(userId, authHeader)
	}()
	go func() {
		defer wg.Done()
		bookings, _ = functions.FetchUserBookingsFromDB(userId)
	}()
	wg.Wait()

	if profileErr != nil {
		return dto.UserDashboard{}, profileErr
	}

	completed := functions.CountBookingsByStatus(bookings, "completed")
	cancelled := functions.CountBookingsByStatus(bookings, "cancelled")
	upcoming := functions.CountBookingsByStatus(bookings, "confirmed")

	recent := bookings
	if len(bookings) > 5 {
		recent = bookings[len(bookings)-5:]
	}

	return dto.UserDashboard{
		UserId:          profile.UserId,
		FullName:        profile.FullName,
		Email:           profile.Email,
		AvatarUrl:       profile.AvatarUrl,
		Role:            profile.Role,
		TotalBookings:   len(bookings),
		UpcomingSlots:   upcoming,
		CompletedVisits: completed,
		CancelledCount:  cancelled,
		RecentBookings:  recent,
	}, nil
}

func FindDashboardApi(c *fiber.Ctx) error {

	role, _ := c.Locals("role").(string)
	userId, _ := c.Locals("userId").(string)
	authHeader := c.Get("Authorization")

	// System-level roles (admin, moderator) always go to the admin dashboard.
	if role == dto.RoleAdmin || role == dto.Rolemoderator {
		dashboard, err := BuildAdminDashboard(authHeader)
		if err != nil {
			return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to build admin dashboard")
		}
		return utils.SendDataResponse(c, dashboard)
	}

	// For everyone else derive the dashboard from ACTUAL ownership records,
	// not from the JWT role claim — a user can be owner AND barber.

	// Owner: has at least one salon where OwnerId == userId.
	salons, _ := functions.FetchUserSalons(userId)
	if len(salons) > 0 {
		dashboard, err := BuildSalonOwnerDashboard(userId, authHeader)
		if err != nil {
			return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to build salon owner dashboard")
		}
		return utils.SendDataResponse(c, dashboard)
	}

	// Barber: has at least one Barber record with UserId == userId.
	barberMemberships, _ := functions.FetchBarberMemberships(userId)
	if len(barberMemberships) > 0 {
		salonId := c.Query("SalonId")
		if salonId == "" {
			salonId = barberMemberships[0].SalonId
		}
		dashboard, err := BuildBarberDashboard(userId, salonId, authHeader)
		if err != nil {
			return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to build barber dashboard")
		}
		return utils.SendDataResponse(c, dashboard)
	}

	// Default: plain customer dashboard.
	dashboard, err := BuildUserDashboard(userId, authHeader)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to build user dashboard")
	}
	return utils.SendDataResponse(c, dashboard)
}
