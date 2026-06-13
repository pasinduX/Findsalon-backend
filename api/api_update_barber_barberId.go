package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateBarberApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)
	barberId := c.Query("BarberId")
	if barberId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "BarberId is required")
	}

	salonId := c.Query("SalonId")
	ownerOk, _ := functions.IsSalonOwner(userId, salonId)
	isOwner := salonId != "" && ownerOk
	barberOk, _ := functions.IsBarberUser(userId)
	isBarber := barberOk

	if !isOwner && !isBarber {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Access denied: must be salon owner or the barber")
	}

	var update bson.M
	if err := c.BodyParser(&update); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	for _, field := range []string{"BarberId", "SalonId", "CreatedAt"} {
		delete(update, field)
	}

	if err := dao.UpdateBarber(barberId, update); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to update barber")
	}
	return utils.SendSuccessResponse(c)
}
