package api

import (
	"log"

	"findsalon-backend/dao"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func DeleteProfileApi(c *fiber.Ctx) error {

	userId := c.Query("UserId")
	if userId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "UserId is required")
	}
	authHeader := c.Get("Authorization")
	if err := dao.DeleteProfile(userId, authHeader); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to delete profile")
	}

	roles, err := dao.FindUserRolesByUserId(userId)
	if err == nil {
		for _, role := range roles {
			if err := dao.DeleteUserRole(role.RoleId); err != nil {
				log.Printf("failed to delete role %s for user %s: %v", role.RoleId, userId, err)
			}
		}
	}

	return utils.SendSuccessResponse(c)
}
