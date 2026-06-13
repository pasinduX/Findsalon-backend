package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginApi(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if req.Email == "" || req.Password == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Email and password are required")
	}

	user, err := dao.FindUserByEmail(req.Email)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password")
	}
	if !user.IsActive {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "Account is inactive")
	}
	if !functions.CheckPasswordHash(req.Password, user.PasswordHash) {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password")
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
			UserId:          user.UserId,
			FullName:        user.FullName,
			Email:           user.Email,
			Phone:           user.Phone,
			AvatarUrl:       user.AvatarUrl,
			GoogleAvatarUrl: user.GoogleAvatarUrl,
			Provider:        user.Provider,
			Role:            user.Role,
			IsActive:        user.IsActive,
			CreatedAt:       user.CreatedAt,
			UpdatedAt:       user.UpdatedAt,
		},
	})
}
