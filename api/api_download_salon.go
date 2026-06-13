package api

import (
	"encoding/json"

	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func DownloadSalonApi(c *fiber.Ctx) error {
	role, _ := c.Locals("role").(string)
	if role != "admin" && role != "moderator" {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Access denied: admin or moderator role required")
	}

	salons, err := dao.FindAllSalons(bson.M{})
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve salons")
	}

	data, err := json.MarshalIndent(salons, "", "  ")
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to marshal salon data")
	}

	c.Set("Content-Disposition", "attachment; filename=salons_export.json")
	c.Set("Content-Type", "application/json")
	return c.Send(data)
}
