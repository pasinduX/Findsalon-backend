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

type BookingEnriched struct {
	BookingId    string    `json:"BookingId"`
	UserId       string    `json:"UserId"`
	SalonId      string    `json:"SalonId"`
	SalonName    string    `json:"SalonName"`
	SalonAddress string    `json:"SalonAddress"`
	SalonArea    string    `json:"SalonArea"`
	BarberId     string    `json:"BarberId"`
	BarberName   string    `json:"BarberName"`
	SlotId       string    `json:"SlotId"`
	ServiceId    string    `json:"ServiceId"`
	Status       string    `json:"Status"`
	BookingType  string    `json:"BookingType"`
	CustomerName string    `json:"CustomerName"`
	Notes        string    `json:"Notes"`
	StartTime    time.Time `json:"StartTime"`
	EndTime      time.Time `json:"EndTime"`
	CreatedAt    time.Time `json:"CreatedAt"`
}

type Booking struct {
	BookingId     string    `json:"BookingId" bson:"BookingId"`
	UserId        string    `json:"UserId" bson:"UserId"`
	SalonId       string    `json:"SalonId" bson:"SalonId" validate:"required"`
	BarberId      string    `json:"BarberId" bson:"BarberId" validate:"required"`
	SlotId        string    `json:"SlotId" bson:"SlotId"`
	ServiceId     string    `json:"ServiceId" bson:"ServiceId"`
	Status        string    `json:"Status" bson:"Status"`
	BookingType   string    `json:"BookingType" bson:"BookingType"`
	CustomerName  string    `json:"CustomerName" bson:"CustomerName"`
	CustomerPhone string    `json:"CustomerPhone" bson:"CustomerPhone"`
	Notes         string    `json:"Notes" bson:"Notes"`
	StartTime     time.Time `json:"StartTime" bson:"StartTime"`
	EndTime       time.Time `json:"EndTime" bson:"EndTime"`
	CreatedAt     time.Time `json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt     time.Time `json:"UpdatedAt" bson:"UpdatedAt"`
	Deleted       bool      `json:"Deleted" bson:"Deleted"`
}
