package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func FindUserApi(c *fiber.Ctx) error {
	userId := c.Query("UserId")
	if userId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "UserId query parameter is required")
	}

	user, err := dao.FindUserByUserId(userId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	response := dto.UserResponse{
		UserId:          user.UserId,
		FullName:        user.FullName,
		Email:           user.Email,
		Phone:           user.Phone,
		AvatarUrl:       user.AvatarUrl,
		GoogleAvatarUrl: user.GoogleAvatarUrl,
		Provider:        user.Provider,
		GoogleId:        user.GoogleId,
		Role:            user.Role,
		IsActive:        user.IsActive,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}

	return utils.SendDataResponse(c, response)
}
