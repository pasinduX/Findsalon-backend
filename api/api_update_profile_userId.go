package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func UpdateProfileApi(c *fiber.Ctx) error {

	userId := c.Query("UserId")
	if userId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "UserId is required")
	}
	currentUserId, _ := c.Locals("userId").(string)
	role, _ := c.Locals("role").(string)
	if userId != currentUserId && role != "admin" {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "cannot update other users profiles")
	}

	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}
	protectedFields := []string{"Email", "Provider", "GoogleId", "GoogleAvatarUrl", "Role", "IsActive", "Deleted", "CreatedAt", "UserId", "PasswordHash", "RefreshToken"}
	for _, field := range protectedFields {
		delete(updateData, field)
	}

	authHeader := c.Get("Authorization")
	if err := dao.UpdateProfile(userId, updateData, authHeader); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to update profile")
	}
	return utils.SendSuccessResponse(c)
}
