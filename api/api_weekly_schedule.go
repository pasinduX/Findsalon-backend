package api

import (
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// POST /api/v1/booking/CreateWeeklySchedule
// Saves (or replaces) the barber's weekly schedule and generates slots for the next 30 days.
func CreateWeeklyScheduleApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)

	var body dto.WeeklySchedule
	if err := c.BodyParser(&body); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	if body.SlotDuration <= 0 {
		body.SlotDuration = 30
	}

	isAllowed, err := functions.IsSalonOwnerOrBarber(userId, body.SalonId)
	if err != nil || !isAllowed {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "You are not authorized")
	}

	if body.ScheduleId == "" {
		body.ScheduleId = uuid.New().String()
	}
	body.UpdatedAt = time.Now()
	body.CreatedAt = time.Now()

	if err := dao.UpsertWeeklySchedule(body); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to save schedule")
	}

	// Generate slots for the next 30 days
	generated := 0
	today := time.Now()
	for i := 0; i < 30; i++ {
		date := today.AddDate(0, 0, i).Format("2006-01-02")
		dayName := today.AddDate(0, 0, i).Weekday().String()

		var workDay *dto.WorkDay
		for _, wd := range body.WorkDays {
			if wd.Day == dayName && wd.IsWorkDay {
				workDay = &wd
				break
			}
		}
		if workDay == nil {
			continue
		}

		blocks, _ := dao.FindScheduleBlocksByBarberDate(body.BarberId, date)
		_ = dao.DeleteNonBookedSlotsForDateBarber(body.BarberId, date)

		slots, err := functions.GenerateTimeSlotsForDate(
			body.BarberId, body.SalonId, date,
			workDay.StartTime, workDay.EndTime,
			body.SlotDuration, blocks,
		)
		if err != nil || len(slots) == 0 {
			continue
		}
		_ = dao.BulkCreateTimeSlots(slots)
		generated += len(slots)
	}

	return utils.SendDataResponse(c, fiber.Map{
		"Schedule":       body,
		"SlotsGenerated": generated,
	})
}

// GET /api/v1/booking/GetWeeklySchedule?BarberId=
func GetWeeklyScheduleApi(c *fiber.Ctx) error {
	barberId := c.Query("BarberId")
	if barberId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "BarberId is required")
	}
	schedule, err := dao.FindWeeklyScheduleByBarber(barberId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Schedule not found")
	}
	return utils.SendDataResponse(c, schedule)
}
