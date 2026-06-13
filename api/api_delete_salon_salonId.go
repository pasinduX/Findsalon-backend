package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func DeleteSalonApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)
	salonId := c.Query("SalonId")
	if salonId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SalonId is required")
	}

	if isOwner, _ := functions.IsSalonOwner(userId, salonId); !isOwner {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Access denied: you are not the owner of this salon")
	}

	if err := dao.DeleteSalon(salonId); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete salon")
	}
	return utils.SendSuccessResponse(c)
}
