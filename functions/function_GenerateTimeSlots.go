package functions

import (
	"fmt"
	"time"

	"findsalon-backend/dto"

	"github.com/google/uuid"
)

const timeLayout = "15:04"

func GenerateTimeSlotsForDate(
	barberId, salonId, date string,
	workStart, workEnd string,
	slotDuration int,
	blocks []dto.ScheduleBlock,
) ([]dto.TimeSlot, error) {
	start, err := time.Parse(timeLayout, workStart)
	if err != nil {
		return nil, fmt.Errorf("invalid workStart: %w", err)
	}
	end, err := time.Parse(timeLayout, workEnd)
	if err != nil {
		return nil, fmt.Errorf("invalid workEnd: %w", err)
	}

	type parsedBlock struct {
		start time.Time
		end   time.Time
		id    string
	}
	var pb []parsedBlock
	for _, b := range blocks {
		bs, err := time.Parse(timeLayout, b.StartTime.UTC().Format(timeLayout))
		if err != nil {
			continue
		}
		be, err := time.Parse(timeLayout, b.EndTime.UTC().Format(timeLayout))
		if err != nil {
			continue
		}
		pb = append(pb, parsedBlock{bs, be, b.BlockId})
	}

	dur := time.Duration(slotDuration) * time.Minute
	now := time.Now()
	var slots []dto.TimeSlot

	for cursor := start; !cursor.Add(dur).After(end); cursor = cursor.Add(dur) {
		slotEnd := cursor.Add(dur)
		status := dto.SlotStatusAvailable
		blockId := ""
		for _, b := range pb {
			if cursor.Before(b.end) && slotEnd.After(b.start) {
				status = dto.SlotStatusBlocked
				blockId = b.id
				break
			}
		}
		slots = append(slots, dto.TimeSlot{
			SlotId:    uuid.New().String(),
			BarberId:  barberId,
			SalonId:   salonId,
			Date:      date,
			StartTime: cursor.Format(timeLayout),
			EndTime:   slotEnd.Format(timeLayout),
			Status:    status,
			IsBooked:  false,
			BlockId:   blockId,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	return slots, nil
}

func DayNameFromDate(date string) (string, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", err
	}
	return t.Weekday().String(), nil
}
