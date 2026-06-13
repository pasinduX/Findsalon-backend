package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func FindAllSpecialtyApi(c *fiber.Ctx) error {
	specialties, err := dao.FindAllSpecialties()
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch specialties")
	}
	return utils.SendDataResponse(c, specialties)
}
