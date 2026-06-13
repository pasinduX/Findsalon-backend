package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

func RegisterApi(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Email, password, and full name are required")
	}

	if !functions.UniqueCheck(dbConfig.USERS_COLLECTION, "Email", req.Email) {
		return utils.SendErrorResponse(c, fiber.StatusConflict, "Email already exists")
	}

	hashed, err := functions.HashPassword(req.Password)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to process password")
	}

	now := time.Now()
	user := dto.User{
		UserId:       uuid.New().String(),
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: hashed,
		Role:         "user",
		Provider:     "local",
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := dao.CreateUser(user); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user")
	}

	accessToken, err := functions.GenerateJWT(user.UserId, user.Email, user.Role)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}
	refreshToken, err := functions.GenerateRefreshToken(user.UserId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate refresh token")
	}

	_ = dao.UpdateUser(user.UserId, bson.M{"RefreshToken": refreshToken})

	return utils.SendDataResponse(c, fiber.Map{
		"AccessToken":  accessToken,
		"RefreshToken": refreshToken,
		"User": dto.UserResponse{
			UserId:    user.UserId,
			FullName:  user.FullName,
			Email:     user.Email,
			Role:      user.Role,
			Provider:  user.Provider,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}
