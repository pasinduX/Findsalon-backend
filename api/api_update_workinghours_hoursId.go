package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateWorkingHoursApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)
	hoursId := c.Query("HoursId")
	salonId := c.Query("SalonId")

	if hoursId == "" || salonId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "HoursId and SalonId are required")
	}

	if isOwner, _ := functions.IsSalonOwner(userId, salonId); !isOwner {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Access denied: you are not the owner of this salon")
	}

	var update bson.M
	if err := c.BodyParser(&update); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	for _, field := range []string{"HoursId", "SalonId", "CreatedAt", "Deleted"} {
		delete(update, field)
	}

	if err := dao.UpdateWorkingHours(hoursId, update); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to update working hours")
	}
	return utils.SendSuccessResponse(c)
}
