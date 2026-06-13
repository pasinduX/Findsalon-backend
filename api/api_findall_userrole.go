package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func FindAllUserRoleApi(c *fiber.Ctx) error {

	pagination := functions.GetPaginationParams(c)
	roleFilter := c.Query("Role")
	var roles interface{}
	var total int
	var err error
	if roleFilter == "" {
		roles, total, err = dao.FindAllUserRoles(int64((pagination.Page-1)*pagination.PageSize), int64(pagination.PageSize))
	} else {
		roles, total, err = dao.FindAllUserRolesByRole(roleFilter, int64((pagination.Page-1)*pagination.PageSize), int64(pagination.PageSize))
	}
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to fetch user roles")
	}

	response := functions.BuildPaginatedResponse(roles, total, pagination)
	return utils.SendDataResponse(c, response)
}
