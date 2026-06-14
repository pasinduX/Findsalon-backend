package api

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/integrations"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type sendPhoneOtpRequest struct {
	UserId string `json:"UserId"`
	Phone  string `json:"Phone"`
}

func SendPhoneOtpApi(c *fiber.Ctx) error {
	var req sendPhoneOtpRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if req.UserId == "" || req.Phone == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "UserId and Phone are required")
	}

	user, err := dao.FindUserByUserId(req.UserId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "User not found")
	}
	if user.PhoneVerified {
		return utils.SendErrorResponse(c, fiber.StatusConflict, "Phone number is already verified")
	}

	otp, err := generateOTP()
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate OTP")
	}
	expiry := time.Now().Add(10 * time.Minute)

	if err := dao.UpdateUser(req.UserId, bson.M{
		"Phone":         req.Phone,
		"PhoneOTP":      otp,
		"PhoneOTPExpiry": expiry,
	}); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to store OTP")
	}

	message := fmt.Sprintf("Your FindSalon verification code is: *%s*\nValid for 10 minutes. Do not share this code.", otp)
	if sendErr := functions.SendWhatsAppMessage(req.Phone, message); sendErr != nil && integrations.TwilioWhatsAppEnabled {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to send WhatsApp message: "+sendErr.Error())
	}

	resp := fiber.Map{
		"message": "Verification code sent to your WhatsApp",
		"phone":   req.Phone,
	}
	// Return OTP in response when Twilio is disabled (development mode)
	if !integrations.TwilioWhatsAppEnabled {
		resp["dev_otp"] = otp
	}

	return utils.SendDataResponse(c, resp)
}

func generateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
