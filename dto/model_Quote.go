package dto

import "time"

type Quote struct {
	QuoteId   string    `json:"QuoteId" bson:"QuoteId"`
	SalonId   string    `json:"SalonId" bson:"SalonId" validate:"required"`
	Text      string    `json:"Text" bson:"Text" validate:"required"`
	Author    string    `json:"Author" bson:"Author"`
	Role      string    `json:"Role" bson:"Role"`
	CreatedAt time.Time `json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt" bson:"UpdatedAt"`
	Deleted   bool      `json:"Deleted" bson:"Deleted"`
}
