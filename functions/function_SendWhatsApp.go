package functions

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"findsalon-backend/integrations"
)

func SendWhatsAppMessage(to, message string) error {
	if !integrations.TwilioWhatsAppEnabled {
		return nil
	}

	accountSid := strings.TrimSpace(integrations.TwilioAccountSid)
	authToken := strings.TrimSpace(integrations.TwilioAuthToken)
	from := normalizeWhatsAppNumber(integrations.TwilioWhatsAppFrom)
	to = normalizeWhatsAppNumber(to)
	message = strings.TrimSpace(message)

	if accountSid == "" || authToken == "" {
		return fmt.Errorf("Twilio credentials not configured")
	}
	if from == "" {
		return fmt.Errorf("Twilio WhatsApp sender not configured")
	}
	if to == "" {
		return fmt.Errorf("WhatsApp recipient is required")
	}
	if message == "" {
		return fmt.Errorf("WhatsApp message is required")
	}

	form := url.Values{}
	form.Set("To", to)
	form.Set("From", from)
	form.Set("Body", message)

	endpoint := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("Twilio WhatsApp request failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

func normalizeWhatsAppNumber(number string) string {
	number = strings.TrimSpace(number)
	if number == "" {
		return ""
	}
	if strings.HasPrefix(number, "whatsapp:") {
		return number
	}
	return "whatsapp:" + number
}
