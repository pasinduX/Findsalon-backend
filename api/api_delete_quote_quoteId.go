package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func DeleteQuoteApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)
	quoteId := c.Query("QuoteId")
	salonId := c.Query("SalonId")

	if quoteId == "" || salonId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "QuoteId and SalonId are required")
	}

	if isOwner, _ := functions.IsSalonOwner(userId, salonId); !isOwner {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Access denied: you are not the owner of this salon")
	}

	if err := dao.DeleteQuote(quoteId); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete quote")
	}
	return utils.SendSuccessResponse(c)
}
