package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func UploadAvatarApi(c *fiber.Ctx) error {

	userId := c.Query("UserId")
	if userId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "UserId is required")
	}
	currentUserId, _ := c.Locals("userId").(string)
	role, _ := c.Locals("role").(string)
	if userId != currentUserId && role != "admin" {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "cannot upload avatar for another user")
	}

	imageUrl, err := functions.SaveUploadedImage(c, "avatar", "avatars/"+userId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	authHeader := c.Get("Authorization")
	if err := dao.UpdateProfile(userId, map[string]interface{}{"AvatarUrl": imageUrl}, authHeader); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to update avatar url")
	}

	return utils.SendDataResponse(c, fiber.Map{"AvatarUrl": imageUrl})
}
