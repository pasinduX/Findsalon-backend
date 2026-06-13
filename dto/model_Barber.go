package dto

import "time"

type Barber struct {
	BarberId   string    `json:"BarberId" bson:"BarberId"`
	SalonId    string    `json:"SalonId" bson:"SalonId" validate:"required"`
	UserId     string    `json:"UserId" bson:"UserId"`
	Name       string    `json:"Name" bson:"Name" validate:"required"`
	Specialties []string  `json:"Specialties" bson:"Specialties"`
	Bio        string    `json:"Bio" bson:"Bio"`
	ImageUrl   string    `json:"ImageUrl" bson:"ImageUrl"`
	ServiceIds []string  `json:"ServiceIds" bson:"ServiceIds"`
	IsActive   bool      `json:"IsActive" bson:"IsActive"`
	CreatedAt  time.Time `json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt  time.Time `json:"UpdatedAt" bson:"UpdatedAt"`
	Deleted    bool      `json:"Deleted" bson:"Deleted"`
}
