package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func DeleteUserRoleApi(c *fiber.Ctx) error {

	roleId := c.Query("RoleId")
	if roleId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "RoleId is required")
	}

	roles, _, err := dao.FindAllUserRoles(0, 10000)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to locate role record")
	}
	var foundRole *dto.UserRole
	for _, role := range roles {
		if role.RoleId == roleId {
			foundRole = &role
			break
		}
	}
	if foundRole == nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "role record not found")
	}

	currentUserId, _ := c.Locals("userId").(string)
	if foundRole.UserId == currentUserId && foundRole.Role == dto.RoleAdmin {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Cannot remove your own admin role")
	}

	if err := dao.DeleteUserRole(roleId); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to delete user role")
	}

	return utils.SendSuccessResponse(c)
}
