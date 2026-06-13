package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func DeleteTemplateApi(c *fiber.Ctx) error {
	templateId := c.Query("TemplateId")
	if templateId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "TemplateId is required")
	}

	_, err := dao.FindTemplateById(templateId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Template not found")
	}

	if err := dao.DeleteTemplate(templateId); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete template")
	}
	return utils.SendSuccessResponse(c)
}
