package dto

import "time"

// SalonService is the unified service model used by both the salon management
// (CRUD) and the booking engine (slot duration + buffer calculation).
type SalonService struct {
	ServiceId   string    `json:"ServiceId" bson:"ServiceId"`
	SalonId     string    `json:"SalonId" bson:"SalonId" validate:"required"`
	Name        string    `json:"Name" bson:"Name" validate:"required"`
	Description string    `json:"Description" bson:"Description"`
	Price       int       `json:"Price" bson:"Price"`
	DurationMin int       `json:"DurationMin" bson:"DurationMin"`
	// BufferMin is cleanup/transition time added after DurationMin.
	// TotalDuration = DurationMin + BufferMin is the window blocked in the calendar.
	BufferMin   int       `json:"BufferMin" bson:"BufferMin"`
	SortOrder   int       `json:"SortOrder" bson:"SortOrder"`
	CreatedAt   time.Time `json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt" bson:"UpdatedAt"`
	Deleted     bool      `json:"Deleted" bson:"Deleted"`
}

// TotalDuration returns the full calendar window this service occupies.
func (s SalonService) TotalDuration() time.Duration {
	return time.Duration(s.DurationMin+s.BufferMin) * time.Minute
}
