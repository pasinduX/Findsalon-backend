package api

import (
	"sync"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func BuildAdminDashboard(authHeader string) (dto.AdminDashboard, error) {
	var totalSalons, activeSalons, totalBarbers, usersTotal int
	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		defer wg.Done()
		totalSalons, activeSalons = functions.FetchAllSalonsCount()
	}()
	go func() {
		defer wg.Done()
		totalBarbers = functions.FetchAllBarbersCount()
	}()
	go func() {
		defer wg.Done()
		_, usersTotal, _ = dao.FindAllProfiles(authHeader, 0, 1)
	}()
	wg.Wait()

	// Total bookings count: fetch summaries and count them.
	bookings, _ := functions.FetchSalonBookingsFromDB("")
	totalBookings := len(bookings)

	return dto.AdminDashboard{
		TotalUsers:    usersTotal,
		TotalSalons:   totalSalons,
		ActiveSalons:  activeSalons,
		TotalBarbers:  totalBarbers,
		TotalBookings: totalBookings,
	}, nil
}

func FindAdminDashboardApi(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	dashboard, err := BuildAdminDashboard(authHeader)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to build admin dashboard")
	}
	return utils.SendDataResponse(c, dashboard)
}
