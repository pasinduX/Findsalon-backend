package functions

import (
	"context"
	"log"
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingNotificationPayload struct {
	BookingId     string
	UserId        string
	SalonId       string
	BarberId      string
	CustomerName  string
	CustomerEmail string
	Date          string
	StartTime     string
	EndTime       string
	EventType     string
}

// NotifyBooking creates an in-app notification and sends an email (if SMTP is enabled).
// Replaces the HTTP call to notification-service in the old booking-ms.
// Designed to be called in a goroutine so booking creation is not blocked.
func NotifyBooking(payload BookingNotificationPayload) {
	title, body := buildNotificationText(payload)

	notification := dto.Notification{
		NotificationId: uuid.New().String(),
		UserId:         payload.UserId,
		Title:          title,
		Body:           body,
		Type:           dto.NotificationTypeBooking,
		EventType:      payload.EventType,
		RefId:          payload.BookingId,
		IsRead:         false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := dao.CreateNotification(notification); err != nil {
		log.Printf("NotifyBooking: failed to save notification: %v", err)
	}

	// Send email via template if SMTP is configured
	if payload.CustomerEmail != "" {
		go sendBookingEmail(payload, title, body)
	}
}

func buildNotificationText(p BookingNotificationPayload) (title, body string) {
	switch p.EventType {
	case dto.EventBookingCreated:
		title = "Booking Confirmed"
		body = "Your booking has been confirmed for " + p.Date + " at " + p.StartTime
	case dto.EventBookingCancelled:
		title = "Booking Cancelled"
		body = "Your booking on " + p.Date + " at " + p.StartTime + " has been cancelled"
	case dto.EventBookingCompleted:
		title = "Visit Completed"
		body = "Thank you for your visit on " + p.Date + ". Please leave a review!"
	default:
		title = "FindSalon Notification"
		body = "You have a new notification regarding booking " + p.BookingId
	}
	return
}

func sendBookingEmail(payload BookingNotificationPayload, title, body string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	col := dbConfig.DATABASE.Collection(dbConfig.TEMPLATES_COLLECTION)
	var tmpl dto.Template
	err := col.FindOne(ctx, bson.M{"EventType": payload.EventType, "IsActive": true, "Deleted": false}).Decode(&tmpl)
	if err != nil {
		// No template configured — send plain body
		_ = SendEmail(payload.CustomerEmail, title, body)
		return
	}

	data := dto.TemplateData{
		CustomerName: payload.CustomerName,
		Date:         payload.Date,
		StartTime:    payload.StartTime,
		EndTime:      payload.EndTime,
		BookingId:    payload.BookingId,
	}
	rendered, err := RenderTemplate(tmpl.BodyTemplate, data)
	if err != nil {
		_ = SendEmail(payload.CustomerEmail, tmpl.Subject, body)
		return
	}
	_ = SendEmail(payload.CustomerEmail, tmpl.Subject, rendered)
}
