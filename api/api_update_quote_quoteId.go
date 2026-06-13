package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateQuoteApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)
	quoteId := c.Query("QuoteId")
	salonId := c.Query("SalonId")

	if quoteId == "" || salonId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "QuoteId and SalonId are required")
	}

	if isOwner, _ := functions.IsSalonOwner(userId, salonId); !isOwner {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Access denied: you are not the owner of this salon")
	}

	var update bson.M
	if err := c.BodyParser(&update); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	for _, field := range []string{"QuoteId", "SalonId", "CreatedAt", "Deleted"} {
		delete(update, field)
	}

	if err := dao.UpdateQuote(quoteId, update); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to update quote")
	}
	return utils.SendSuccessResponse(c)
}
