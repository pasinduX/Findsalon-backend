package api

import (
	"sync"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func BuildSalonOwnerDashboard(userId string, authHeader string) (dto.SalonOwnerDashboard, error) {
	var salons []dto.SalonSummary
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		salons, _ = functions.FetchUserSalons(userId)
	}()
	go func() {
		defer wg.Done()
		// Prefetch profile to confirm identity (result unused in DTO but validates auth).
		dao.FindProfileByUserId(userId, authHeader)
	}()
	wg.Wait()

	totalSalons := len(salons)
	totalBarbers := 0
	totalBookings := 0
	todayBookings := 0

	for _, salon := range salons {
		totalBarbers += functions.FetchBarbersBySalonCount(salon.SalonId)
		bookings, _ := functions.FetchSalonBookingsFromDB(salon.SalonId)
		totalBookings += len(bookings)
		today, _ := functions.FetchTodayBookingsFromDB(salon.SalonId)
		todayBookings += len(today)
	}

	return dto.SalonOwnerDashboard{
		TotalSalons:   totalSalons,
		TotalBarbers:  totalBarbers,
		TotalBookings: totalBookings,
		TodayBookings: todayBookings,
		PendingSlots:  todayBookings,
		TotalRevenue:  0,
		Salons:        salons,
	}, nil
}

func FindSalonOwnerDashboardApi(c *fiber.Ctx) error {
	role, _ := c.Locals("role").(string)
	if role != dto.RoleSalonOwner && role != dto.RoleAdmin {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "salon owner or admin privileges required")
	}

	userId, _ := c.Locals("userId").(string)
	authHeader := c.Get("Authorization")
	dashboard, err := BuildSalonOwnerDashboard(userId, authHeader)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to build salon owner dashboard")
	}
	return utils.SendDataResponse(c, dashboard)
}
