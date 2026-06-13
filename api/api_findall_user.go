package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func FindAllUserApi(c *fiber.Ctx) error {
	users, err := dao.FindAllUsers()
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve users")
	}

	responses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, dto.UserResponse{
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
		})
	}

	return utils.SendDataResponse(c, responses)
}
