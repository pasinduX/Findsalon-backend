package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func FindTimeSlotApi(c *fiber.Ctx) error {
	slotId := c.Query("SlotId")
	if slotId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SlotId is required")
	}

	slot, err := dao.FindTimeSlotBySlotId(slotId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Time slot not found")
	}

	return utils.SendDataResponse(c, slot)
}
