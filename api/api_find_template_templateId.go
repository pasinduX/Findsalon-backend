package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func FindTemplateApi(c *fiber.Ctx) error {
	templateId := c.Query("TemplateId")
	if templateId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "TemplateId is required")
	}

	template, err := dao.FindTemplateById(templateId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Template not found")
	}

	return utils.SendDataResponse(c, template)
}
