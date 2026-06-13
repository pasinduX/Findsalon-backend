package dto

import "time"

const (
	NotificationTypeBooking = "booking"
	NotificationTypeSystem  = "system"
	NotificationTypePromo   = "promo"
	NotificationTypeReview  = "review"
)

const (
	EventBookingCreated   = "booking_created"
	EventBookingCancelled = "booking_cancelled"
	EventBookingCompleted = "booking_completed"
	EventReviewReceived   = "review_received"
	EventCustom           = "custom"
	EventBulk             = "bulk"
)

type Notification struct {
	NotificationId string    `json:"NotificationId" bson:"NotificationId"`
	UserId         string    `json:"UserId" bson:"UserId"`
	Title          string    `json:"Title" bson:"Title"`
	Body           string    `json:"Body" bson:"Body"`
	Type           string    `json:"Type" bson:"Type"`
	EventType      string    `json:"EventType" bson:"EventType"`
	RefId          string    `json:"RefId" bson:"RefId"`
	IsRead         bool      `json:"IsRead" bson:"IsRead"`
	CreatedAt      time.Time `json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt      time.Time `json:"UpdatedAt" bson:"UpdatedAt"`
	Deleted        bool      `json:"Deleted" bson:"Deleted"`
}
