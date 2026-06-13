package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/integrations"
	"findsalon-backend/utils"
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func GoogleCallbackApi(c *fiber.Ctx) error {
	state := c.Query("state")
	cookieState := c.Cookies("oauth_state")

	if state == "" || state != cookieState {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid OAuth state")
	}

	c.Cookie(&fiber.Cookie{
		Name:   "oauth_state",
		Value:  "",
		MaxAge: -1,
	})

	code := c.Query("code")
	if code == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Authorization code is required")
	}

	token, err := functions.ExchangeGoogleCode(code)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to exchange authorization code")
	}

	googleUser, err := functions.GetGoogleUserInfo(token)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to get Google user info")
	}

	var user dto.User
	existingUser, err := dao.FindUserByEmail(googleUser.Email)
	if err != nil {
		if err.Error() != "user not found" {
			return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Database error")
		}

		newUser := dto.User{
			UserId:          uuid.New().String(),
			FullName:        googleUser.Name,
			Email:           googleUser.Email,
			AvatarUrl:       googleUser.Picture,
			GoogleAvatarUrl: googleUser.Picture,
			GoogleId:        googleUser.Id,
			Provider:        "google",
			Role:            "user",
			IsActive:        true,
			Deleted:         false,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		if createErr := dao.CreateUser(newUser); createErr != nil {
			return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user")
		}
		user = newUser
	} else {
		user = existingUser
	}

	accessToken, err := functions.GenerateJWT(user.UserId, user.Email, user.Role)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate access token")
	}

	refreshToken, err := functions.GenerateRefreshToken(user.UserId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate refresh token")
	}

	if err := dao.UpdateUser(user.UserId, bson.M{"RefreshToken": refreshToken}); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to save refresh token")
	}

	frontendURL := integrations.FrontendUrl
	if frontendURL == "" {
		frontendURL = "http://localhost:4000"
	}

	redirectTo := fmt.Sprintf(
		"%s/auth/callback?access_token=%s&refresh_token=%s",
		frontendURL,
		url.QueryEscape(accessToken),
		url.QueryEscape(refreshToken),
	)
	return c.Redirect(redirectTo, fiber.StatusFound)
}
