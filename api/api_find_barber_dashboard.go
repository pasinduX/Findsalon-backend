package api

import (
	"fmt"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func BuildBarberDashboard(userId string, salonId string, authHeader string) (dto.BarberDashboard, error) {
	if salonId == "" {
		return dto.BarberDashboard{}, fmt.Errorf("SalonId is required")
	}

	barber, err := dao.FindBarberByUserId(userId)
	if err != nil || barber.BarberId == "" {
		return dto.BarberDashboard{}, fmt.Errorf("Barber profile not found for this user")
	}

	salonSummary, _ := functions.FetchSalonSummaryById(salonId)
	memberships, _ := functions.FetchBarberMemberships(userId)

	bookings, _ := functions.FetchUserBookingsFromDB(barber.BarberId)
	todayBookings := functions.FetchBarbersBySalonCount(salonId)
	completed := functions.CountBookingsByStatus(bookings, "completed")

	// Average rating from reviews.
	ratingSummary, _ := dao.AggregateRatingSummary("BarberId", barber.BarberId)
	avgRating := ratingSummary.AverageRating

	return dto.BarberDashboard{
		SalonId:        salonId,
		SalonName:      salonSummary.Name,
		Memberships:    memberships,
		TodaySlots:     todayBookings,
		TodayBookings:  todayBookings,
		TotalCompleted: completed,
		AverageRating:  avgRating,
	}, nil
}

func FindBarberDashboardApi(c *fiber.Ctx) error {

	userId, _ := c.Locals("userId").(string)
	authHeader := c.Get("Authorization")
	salonId := c.Query("SalonId")

	if salonId == "" {
		memberships, _ := functions.FetchBarberMemberships(userId)
		if len(memberships) == 0 {
			return utils.SendErrorResponse(c, fiber.StatusForbidden, "No barber profile found for this account")
		}
		salonId = memberships[0].SalonId
	}

	dashboard, err := BuildBarberDashboard(userId, salonId, authHeader)
	if err != nil {
		if err.Error() == "Barber profile not found for this user" {
			return utils.SendErrorResponse(c, fiber.StatusNotFound, err.Error())
		}
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to build barber dashboard")
	}
	return utils.SendDataResponse(c, dashboard)
}
