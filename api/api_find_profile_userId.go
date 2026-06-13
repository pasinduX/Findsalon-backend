package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func FindProfileApi(c *fiber.Ctx) error {

	requestedUserId := c.Query("UserId")
	currentUserId, _ := c.Locals("userId").(string)
	role, _ := c.Locals("role").(string)
	if requestedUserId == "" {
		requestedUserId = currentUserId
	}
	if requestedUserId != currentUserId && role != "admin" && role != "moderator" {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Cannot view other users profiles")
	}

	authHeader := c.Get("Authorization")
	profile, err := dao.FindProfileByUserId(requestedUserId, authHeader)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to fetch profile")
	}
	return utils.SendDataResponse(c, profile)
}
