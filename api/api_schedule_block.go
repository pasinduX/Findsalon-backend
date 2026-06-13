package api

import (
	"fmt"
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// blockTimeToUTC converts a date ("YYYY-MM-DD") and HH:MM string from the
// barber's schedule timezone to a UTC time.Time.
func blockTimeToUTC(barberId, date, hhmm string) (time.Time, error) {
	tz := "UTC"
	if sched, err := dao.FindWeeklyScheduleByBarber(barberId); err == nil && sched != nil {
		tz = sched.EffectiveTimezone()
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone %s: %w", tz, err)
	}
	localDate, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date %q: %w", date, err)
	}
	t, err := time.Parse("15:04", hhmm)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time %q: %w", hhmm, err)
	}
	return time.Date(
		localDate.Year(), localDate.Month(), localDate.Day(),
		t.Hour(), t.Minute(), 0, 0, loc,
	).UTC(), nil
}

// regenerateSlotsForDate deletes non-booked slots and regenerates them (legacy flow).
func regenerateSlotsForDate(barberId, salonId, date string) error {
	schedule, err := dao.FindWeeklyScheduleByBarber(barberId)
	if err != nil {
		return nil
	}
	dayName, _ := functions.DayNameFromDate(date)
	var workDay *dto.WorkDay
	for _, wd := range schedule.WorkDays {
		if wd.Day == dayName && wd.Active() {
			workDay = &wd
			break
		}
	}
	if workDay == nil {
		return dao.DeleteNonBookedSlotsForDateBarber(barberId, date)
	}
	blocks, _ := dao.FindScheduleBlocksByBarberDate(barberId, date)
	_ = dao.DeleteNonBookedSlotsForDateBarber(barberId, date)
	slots, err := functions.GenerateTimeSlotsForDate(
		barberId, salonId, date,
		workDay.StartTime, workDay.EndTime,
		schedule.EffectiveSlotStep(), blocks,
	)
	if err != nil || len(slots) == 0 {
		return nil
	}
	return dao.BulkCreateTimeSlots(slots)
}

// POST /api/v1/booking/CreateScheduleBlock
func CreateScheduleBlockApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)
	var req dto.CreateScheduleBlockRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	isAllowed, err := functions.IsSalonOwnerOrBarber(userId, req.SalonId)
	if err != nil || !isAllowed {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "You are not authorized")
	}
	startUTC, err := blockTimeToUTC(req.BarberId, req.Date, req.StartTime)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid StartTime: "+err.Error())
	}
	endUTC, err := blockTimeToUTC(req.BarberId, req.Date, req.EndTime)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid EndTime: "+err.Error())
	}
	now := time.Now().UTC()
	block := dto.ScheduleBlock{
		BlockId:   uuid.New().String(),
		BarberId:  req.BarberId,
		SalonId:   req.SalonId,
		StartTime: startUTC,
		EndTime:   endUTC,
		BlockType: req.BlockType,
		Note:      req.Note,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := dao.CreateScheduleBlock(block); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create block")
	}
	_ = regenerateSlotsForDate(block.BarberId, block.SalonId, req.Date)
	return utils.SendDataResponse(c, block)
}

// PUT /api/v1/booking/UpdateScheduleBlock?BlockId=
func UpdateScheduleBlockApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)
	blockId := c.Query("BlockId")
	if blockId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "BlockId is required")
	}
	existing, err := dao.FindScheduleBlockById(blockId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Block not found")
	}
	isAllowed, err := functions.IsSalonOwnerOrBarber(userId, existing.SalonId)
	if err != nil || !isAllowed {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "You are not authorized")
	}
	var req dto.CreateScheduleBlockRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	date := req.Date
	if date == "" {
		date = existing.StartTime.UTC().Format("2006-01-02")
	}
	startUTC, err := blockTimeToUTC(existing.BarberId, date, req.StartTime)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid StartTime: "+err.Error())
	}
	endUTC, err := blockTimeToUTC(existing.BarberId, date, req.EndTime)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid EndTime: "+err.Error())
	}
	updated := dto.ScheduleBlock{
		BlockId:   blockId,
		BarberId:  existing.BarberId,
		SalonId:   existing.SalonId,
		StartTime: startUTC,
		EndTime:   endUTC,
		BlockType: req.BlockType,
		Note:      req.Note,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: time.Now().UTC(),
	}
	if err := dao.UpdateScheduleBlock(blockId, updated); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to update block")
	}
	_ = regenerateSlotsForDate(updated.BarberId, updated.SalonId, date)
	return utils.SendDataResponse(c, updated)
}

// DELETE /api/v1/booking/DeleteScheduleBlock?BlockId=
func DeleteScheduleBlockApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)
	blockId := c.Query("BlockId")
	if blockId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "BlockId is required")
	}
	existing, err := dao.FindScheduleBlockById(blockId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Block not found")
	}
	isAllowed, err := functions.IsSalonOwnerOrBarber(userId, existing.SalonId)
	if err != nil || !isAllowed {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "You are not authorized")
	}
	if err := dao.DeleteScheduleBlock(blockId); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete block")
	}
	date := existing.StartTime.UTC().Format("2006-01-02")
	_ = regenerateSlotsForDate(existing.BarberId, existing.SalonId, date)
	return utils.SendSuccessResponse(c)
}

// GET /api/v1/booking/FindallScheduleBlock?BarberId=&Date=
func FindAllScheduleBlockApi(c *fiber.Ctx) error {
	barberId := c.Query("BarberId")
	if barberId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "BarberId is required")
	}
	date := c.Query("Date")
	blocks, err := dao.FindScheduleBlocksByBarberDate(barberId, date)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch blocks")
	}
	return utils.SendDataResponse(c, blocks)
}
