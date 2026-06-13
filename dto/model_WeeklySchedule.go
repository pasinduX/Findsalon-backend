package dto

import "time"

type WorkDay struct {
	Day       string `json:"Day" bson:"Day"`
	StartTime string `json:"StartTime" bson:"StartTime"`
	EndTime   string `json:"EndTime" bson:"EndTime"`
	IsActive  bool   `json:"IsActive" bson:"IsActive"`
	IsWorkDay bool   `json:"IsWorkDay" bson:"IsWorkDay"`
}

func (w WorkDay) Active() bool { return w.IsActive || w.IsWorkDay }

type WeeklySchedule struct {
	ScheduleId   string    `json:"ScheduleId" bson:"ScheduleId"`
	BarberId     string    `json:"BarberId" bson:"BarberId" validate:"required"`
	SalonId      string    `json:"SalonId" bson:"SalonId" validate:"required"`
	Timezone     string    `json:"Timezone" bson:"Timezone"`
	SlotStep     int       `json:"SlotStep" bson:"SlotStep"`
	SlotDuration int       `json:"SlotDuration" bson:"SlotDuration"`
	WorkDays     []WorkDay `json:"WorkDays" bson:"WorkDays"`
	CreatedAt    time.Time `json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt    time.Time `json:"UpdatedAt" bson:"UpdatedAt"`
}

func (s WeeklySchedule) EffectiveSlotStep() int {
	if s.SlotStep > 0 {
		return s.SlotStep
	}
	if s.SlotDuration > 0 {
		return s.SlotDuration
	}
	return 30
}

func (s WeeklySchedule) EffectiveTimezone() string {
	if s.Timezone != "" {
		return s.Timezone
	}
	return "UTC"
}
