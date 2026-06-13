package dto

import "time"

type WorkingHours struct {
	HoursId   string    `json:"HoursId" bson:"HoursId"`
	SalonId   string    `json:"SalonId" bson:"SalonId" validate:"required"`
	Day       string    `json:"Day" bson:"Day" validate:"required"`
	Hours     string    `json:"Hours" bson:"Hours" validate:"required"`
	SortOrder int       `json:"SortOrder" bson:"SortOrder"`
	CreatedAt time.Time `json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt" bson:"UpdatedAt"`
	Deleted   bool      `json:"Deleted" bson:"Deleted"`
}
