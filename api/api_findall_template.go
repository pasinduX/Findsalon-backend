package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func FindAllTemplateApi(c *fiber.Ctx) error {
	templates, err := dao.FindAllTemplates()
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to load templates")
	}
	return utils.SendDataResponse(c, templates)
}
