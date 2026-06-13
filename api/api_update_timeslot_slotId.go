package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateTimeSlotApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)

	slotId := c.Query("SlotId")
	if slotId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SlotId is required")
	}

	slot, err := dao.FindTimeSlotBySlotId(slotId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Time slot not found")
	}

	isBarber, _ := functions.IsBarberOfSlot(userId, slotId)
	isAllowed, _ := functions.IsSalonOwnerOrBarber(userId, slot.SalonId)
	if !isBarber && !isAllowed {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "You are not authorized to update this time slot")
	}

	if slot.IsBooked {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Cannot update a booked time slot")
	}

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	delete(body, "SlotId")
	delete(body, "BarberId")
	delete(body, "SalonId")
	delete(body, "IsBooked")
	delete(body, "CreatedAt")
	delete(body, "Deleted")

	if len(body) == 0 {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "No valid fields to update")
	}

	if err := dao.UpdateTimeSlot(slotId, bson.M{"$set": body}); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to update time slot")
	}

	return utils.SendSuccessResponse(c)
}
