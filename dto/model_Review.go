package dto

import "time"

// Review is the unified model used by both the booking flow (POST /booking/CreateReview)
// and the standalone review service (/review/ endpoints).
type Review struct {
	ReviewId  string    `json:"ReviewId" bson:"ReviewId"`
	SalonId   string    `json:"SalonId" bson:"SalonId" validate:"required"`
	BarberId  string    `json:"BarberId" bson:"BarberId"`
	BookingId string    `json:"BookingId" bson:"BookingId"`
	UserId    string    `json:"UserId" bson:"UserId" validate:"required"`
	UserName  string    `json:"UserName" bson:"UserName"`
	Rating    int       `json:"Rating" bson:"Rating" validate:"required,min=1,max=5"`
	Comment   string    `json:"Comment" bson:"Comment"`
	IsVisible bool      `json:"IsVisible" bson:"IsVisible"`
	CreatedAt time.Time `json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt" bson:"UpdatedAt"`
	Deleted   bool      `json:"Deleted" bson:"Deleted"`
}

type PaginatedReviewResponse struct {
	Reviews    []Review `json:"Reviews"`
	Page       int      `json:"Page"`
	PageSize   int      `json:"PageSize"`
	TotalCount int      `json:"TotalCount"`
	TotalPages int      `json:"TotalPages"`
}
