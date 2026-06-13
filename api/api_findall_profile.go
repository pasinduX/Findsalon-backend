package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func FindAllProfileApi(c *fiber.Ctx) error {

	pagination := functions.GetPaginationParams(c)
	roleFilter := c.Query("Role")
	authHeader := c.Get("Authorization")
	profiles, total, err := dao.FindAllProfiles(authHeader, (pagination.Page-1)*pagination.PageSize, pagination.PageSize)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to fetch profiles")
	}

	if roleFilter != "" {
		filtered := []dto.Profile{}
		for _, profile := range profiles {
			if profile.Role == roleFilter {
				filtered = append(filtered, profile)
			}
		}
		profiles = filtered
		total = len(filtered)
	}

	response := functions.BuildPaginatedResponse(profiles, total, pagination)
	return utils.SendDataResponse(c, response)
}
