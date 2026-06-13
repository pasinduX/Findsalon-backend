package api

import (
	"strings"

	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateTemplateApi(c *fiber.Ctx) error {
	templateId := c.Query("TemplateId")
	if templateId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "TemplateId is required")
	}

	_, err := dao.FindTemplateById(templateId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Template not found")
	}

	var raw map[string]interface{}
	if err := c.BodyParser(&raw); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	allowed := map[string]bool{
		"Name":         true,
		"Subject":      true,
		"BodyTemplate": true,
		"IsActive":     true,
	}

	update := bson.M{}
	for key, value := range raw {
		if allowed[key] {
			update[key] = value
		}
	}

	if len(update) == 0 {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "No valid fields to update")
	}

	if _, ok := update["IsActive"]; ok {
		if s, ok := update["IsActive"].(string); ok {
			update["IsActive"] = strings.EqualFold(s, "true")
		}
	}

	if err := dao.UpdateTemplate(templateId, update); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to update template")
	}
	return utils.SendSuccessResponse(c)
}
