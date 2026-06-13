package functions

import (
	"bytes"
	"text/template"

	"findsalon-backend/dto"
)

func RenderTemplate(tmplText string, data dto.TemplateData) (string, error) {
	t, err := template.New("email").Parse(tmplText)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func DefaultBookingCreatedTemplate() string {
	return `<h2>Booking Confirmed</h2>
<p>Dear {{.CustomerName}},</p>
<p>Your booking has been confirmed.</p>
<p><strong>Date:</strong> {{.Date}}<br>
<strong>Time:</strong> {{.StartTime}} – {{.EndTime}}<br>
<strong>Barber/Salon:</strong> {{.SalonName}}<br>
<strong>Booking ID:</strong> {{.BookingId}}</p>
<p>Thank you for choosing FindSalon!</p>`
}

func DefaultBookingCancelledTemplate() string {
	return `<h2>Booking Cancelled</h2>
<p>Dear {{.CustomerName}},</p>
<p>Your booking on {{.Date}} at {{.StartTime}} has been cancelled.</p>
<p>Booking ID: {{.BookingId}}</p>
<p>We hope to see you again soon!</p>`
}

func DefaultBookingCompletedTemplate() string {
	return `<h2>Visit Completed</h2>
<p>Dear {{.CustomerName}},</p>
<p>Thank you for your visit on {{.Date}}!</p>
<p>We'd love to hear your feedback. Please leave a review on FindSalon.</p>`
}
