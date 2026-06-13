package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func DeleteSalonServiceApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)
	serviceId := c.Query("ServiceId")
	salonId := c.Query("SalonId")

	if serviceId == "" || salonId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "ServiceId and SalonId are required")
	}

	if isOwner, _ := functions.IsSalonOwner(userId, salonId); !isOwner {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Access denied: you are not the owner of this salon")
	}

	if err := dao.DeleteSalonService(serviceId); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete salon service")
	}
	return utils.SendSuccessResponse(c)
}
