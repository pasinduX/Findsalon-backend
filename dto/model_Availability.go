package dto

import "time"

type AvailabilityRequest struct {
	BarberId           string `query:"BarberId" validate:"required"`
	ServiceId          string `query:"ServiceId" validate:"required"`
	Date               string `query:"Date" validate:"required"`
	IncludeUnavailable bool   `query:"IncludeUnavailable"`
}

type AvailableSlot struct {
	StartTime         time.Time `json:"StartTime"`
	EndTime           time.Time `json:"EndTime"`
	DisplayStart      string    `json:"DisplayStart"`
	DisplayEnd        string    `json:"DisplayEnd"`
	Status            string    `json:"Status"`
	IsAvailable       bool      `json:"IsAvailable"`
	UnavailableReason string    `json:"UnavailableReason,omitempty"`
}

type DirectBookingRequest struct {
	BarberId      string    `json:"BarberId" validate:"required"`
	SalonId       string    `json:"SalonId" validate:"required"`
	ServiceId     string    `json:"ServiceId" validate:"required"`
	CustomerId    string    `json:"CustomerId"`
	CustomerName  string    `json:"CustomerName"`
	CustomerPhone string    `json:"CustomerPhone"`
	StartTime     time.Time `json:"StartTime" validate:"required"`
	Notes         string    `json:"Notes"`
}
