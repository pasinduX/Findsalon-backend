package api

import (
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type verifyPhoneOtpRequest struct {
	UserId string `json:"UserId"`
	Code   string `json:"Code"`
}

func VerifyPhoneOtpApi(c *fiber.Ctx) error {
	var req verifyPhoneOtpRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if req.UserId == "" || req.Code == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "UserId and Code are required")
	}

	user, err := dao.FindUserByUserId(req.UserId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "User not found")
	}
	if user.PhoneVerified {
		return utils.SendErrorResponse(c, fiber.StatusConflict, "Phone number is already verified")
	}
	if user.PhoneOTP == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "No verification code found. Please request a new code.")
	}
	if time.Now().After(user.PhoneOTPExpiry) {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Verification code has expired. Please request a new code.")
	}
	if user.PhoneOTP != req.Code {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid verification code")
	}

	if err := dao.UpdateUser(req.UserId, bson.M{
		"PhoneVerified":  true,
		"PhoneOTP":       "",
		"PhoneOTPExpiry": time.Time{},
	}); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to verify phone")
	}

	return utils.SendDataResponse(c, dto.UserResponse{
		UserId:          user.UserId,
		FullName:        user.FullName,
		Email:           user.Email,
		Phone:           user.Phone,
		PhoneVerified:   true,
		AvatarUrl:       user.AvatarUrl,
		GoogleAvatarUrl: user.GoogleAvatarUrl,
		Provider:        user.Provider,
		GoogleId:        user.GoogleId,
		Role:            user.Role,
		IsActive:        user.IsActive,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	})
}
