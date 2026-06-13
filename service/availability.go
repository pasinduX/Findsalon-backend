package service

import (
	"fmt"
	"time"

	"findsalon-backend/dto"
	salonerr "findsalon-backend/errors"
)

const timeLayout = "15:04"

func overlaps(aStart, aEnd, bStart, bEnd time.Time) bool {
	return aStart.Before(bEnd) && aEnd.After(bStart)
}

type WorkWindow struct {
	Start time.Time
	End   time.Time
}

func GetWorkWindow(schedule dto.WeeklySchedule, dateStr string) (WorkWindow, error) {
	tz := schedule.EffectiveTimezone()
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return WorkWindow{}, fmt.Errorf("%w: %s", salonerr.ErrInvalidTimezone, tz)
	}
	localDate, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		return WorkWindow{}, fmt.Errorf("%w: date must be YYYY-MM-DD", salonerr.ErrInvalidInput)
	}
	weekday := localDate.Weekday().String()
	for _, wd := range schedule.WorkDays {
		if wd.Day != weekday || !wd.Active() {
			continue
		}
		start, err := buildLocalTime(localDate, wd.StartTime, loc)
		if err != nil {
			return WorkWindow{}, fmt.Errorf("invalid StartTime %q: %w", wd.StartTime, err)
		}
		end, err := buildLocalTime(localDate, wd.EndTime, loc)
		if err != nil {
			return WorkWindow{}, fmt.Errorf("invalid EndTime %q: %w", wd.EndTime, err)
		}
		return WorkWindow{Start: start.UTC(), End: end.UTC()}, nil
	}
	return WorkWindow{}, salonerr.ErrNotWorkingDay
}

func buildLocalTime(date time.Time, hhmm string, loc *time.Location) (time.Time, error) {
	t, err := time.Parse(timeLayout, hhmm)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, loc), nil
}

type AvailabilityInput struct {
	Window   WorkWindow
	Service  dto.SalonService
	StepMin  int
	Bookings []dto.Booking
	Blocks   []dto.ScheduleBlock
	Timezone string
}

func GetAvailableSlots(input AvailabilityInput) ([]dto.AvailableSlot, error) {
	loc, err := time.LoadLocation(input.Timezone)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", salonerr.ErrInvalidTimezone, input.Timezone)
	}
	step := time.Duration(input.StepMin) * time.Minute
	duration := input.Service.TotalDuration()
	now := time.Now().UTC()
	var slots []dto.AvailableSlot
	for t := input.Window.Start; ; t = t.Add(step) {
		slotEnd := t.Add(duration)
		if slotEnd.After(input.Window.End) {
			break
		}
		if t.Before(now) {
			continue
		}
		if isAvailable(t, slotEnd, input.Bookings, input.Blocks) {
			local := t.In(loc)
			slots = append(slots, dto.AvailableSlot{
				StartTime:    t,
				EndTime:      slotEnd,
				DisplayStart: local.Format(timeLayout),
				DisplayEnd:   slotEnd.In(loc).Format(timeLayout),
			})
		}
	}
	return slots, nil
}

func isAvailable(start, end time.Time, bookings []dto.Booking, blocks []dto.ScheduleBlock) bool {
	for i := range bookings {
		if bookings[i].Status == dto.BookingStatusCancelled {
			continue
		}
		if overlaps(start, end, bookings[i].StartTime, bookings[i].EndTime) {
			return false
		}
	}
	for i := range blocks {
		if overlaps(start, end, blocks[i].StartTime, blocks[i].EndTime) {
			return false
		}
	}
	return true
}
