package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func DeleteTimeSlotApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)

	slotId := c.Query("SlotId")
	if slotId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SlotId is required")
	}

	slot, err := dao.FindTimeSlotBySlotId(slotId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Time slot not found")
	}

	if slot.IsBooked {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Cannot delete a booked time slot")
	}

	isAllowed, _ := functions.IsSalonOwnerOrBarber(userId, slot.SalonId)
	if !isAllowed {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "You are not authorized to delete this time slot")
	}

	if err := dao.DeleteTimeSlot(slotId); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete time slot")
	}

	return utils.SendSuccessResponse(c)
}
