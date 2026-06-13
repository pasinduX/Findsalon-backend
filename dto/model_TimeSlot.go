package dto

import "time"

type SlotStatus string

const (
	SlotStatusAvailable SlotStatus = "available"
	SlotStatusBooked    SlotStatus = "booked"
	SlotStatusBlocked   SlotStatus = "blocked"
)

type TimeSlot struct {
	SlotId    string     `json:"SlotId" bson:"SlotId"`
	BarberId  string     `json:"BarberId" bson:"BarberId" validate:"required"`
	SalonId   string     `json:"SalonId" bson:"SalonId" validate:"required"`
	Date      string     `json:"Date" bson:"Date" validate:"required"`
	StartTime string     `json:"StartTime" bson:"StartTime" validate:"required"`
	EndTime   string     `json:"EndTime" bson:"EndTime" validate:"required"`
	Status    SlotStatus `json:"Status" bson:"Status"`
	IsBooked  bool       `json:"IsBooked" bson:"IsBooked"`
	BlockId   string     `json:"BlockId" bson:"BlockId"`
	CreatedAt time.Time  `json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt time.Time  `json:"UpdatedAt" bson:"UpdatedAt"`
	Deleted   bool       `json:"Deleted" bson:"Deleted"`
}
