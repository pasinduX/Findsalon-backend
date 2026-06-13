package dto

import "time"

const (
	BookingStatusConfirmed = "confirmed"
	BookingStatusCancelled = "cancelled"
	BookingStatusCompleted = "completed"
	BookingStatusBooked    = "booked"
)

const (
	BookingTypeOnline = "online"
	BookingTypeWalkIn = "walk_in"
	BookingTypeDirect = "direct"
)

type Booking struct {
	BookingId    string    `json:"BookingId" bson:"BookingId"`
	UserId       string    `json:"UserId" bson:"UserId"`
	SalonId      string    `json:"SalonId" bson:"SalonId" validate:"required"`
	BarberId     string    `json:"BarberId" bson:"BarberId" validate:"required"`
	SlotId       string    `json:"SlotId" bson:"SlotId"`
	ServiceId    string    `json:"ServiceId" bson:"ServiceId"`
	Status       string    `json:"Status" bson:"Status"`
	BookingType  string    `json:"BookingType" bson:"BookingType"`
	CustomerName string    `json:"CustomerName" bson:"CustomerName"`
	Notes        string    `json:"Notes" bson:"Notes"`
	StartTime    time.Time `json:"StartTime" bson:"StartTime"`
	EndTime      time.Time `json:"EndTime" bson:"EndTime"`
	CreatedAt    time.Time `json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt    time.Time `json:"UpdatedAt" bson:"UpdatedAt"`
	Deleted      bool      `json:"Deleted" bson:"Deleted"`
}
