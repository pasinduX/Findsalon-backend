package functions

import (
	"findsalon-backend/integrations"

	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, body string) error {
	if !integrations.SmtpEnabled {
		return nil
	}
	m := gomail.NewMessage()
	m.SetHeader("From", integrations.SmtpFromName+" <"+integrations.SmtpFromEmail+">")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(integrations.SmtpHost, integrations.SmtpPort, integrations.SmtpUsername, integrations.SmtpPassword)
	return d.DialAndSend(m)
}

func SendEmailAsync(to, subject, body string) {
	go func() {
		_ = SendEmail(to, subject, body)
	}()
}
