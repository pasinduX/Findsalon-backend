package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type refreshTokenRequest struct {
	RefreshToken string `json:"RefreshToken"`
}

func RefreshTokenApi(c *fiber.Ctx) error {
	var req refreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if req.RefreshToken == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "RefreshToken is required")
	}

	claims, err := functions.ValidateJWT(req.RefreshToken)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired refresh token")
	}

	user, err := dao.FindUserByUserId(claims.UserId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	if user.RefreshToken != req.RefreshToken {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "Refresh token mismatch")
	}

	newAccessToken, err := functions.GenerateJWT(user.UserId, user.Email, user.Role)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate access token")
	}

	newRefreshToken, err := functions.GenerateRefreshToken(user.UserId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate refresh token")
	}

	if err := dao.UpdateUser(user.UserId, bson.M{"RefreshToken": newRefreshToken}); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to update refresh token")
	}

	return utils.SendDataResponse(c, fiber.Map{
		"AccessToken":  newAccessToken,
		"RefreshToken": newRefreshToken,
	})
}
