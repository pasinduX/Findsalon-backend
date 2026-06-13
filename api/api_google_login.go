package api

import (
	"crypto/rand"
	"encoding/base64"
	"findsalon-backend/functions"

	"github.com/gofiber/fiber/v2"
)

func generateStateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func GoogleLoginApi(c *fiber.Ctx) error {
	state := generateStateToken()

	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HTTPOnly: true,
		MaxAge:   600,
		SameSite: "Lax",
	})

	authURL := functions.GetGoogleAuthURL(state)
	return c.Redirect(authURL, fiber.StatusFound)
}
