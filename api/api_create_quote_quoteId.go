package api

import (
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateQuoteApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)

	var quote dto.Quote
	if err := c.BodyParser(&quote); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := validate.Struct(quote); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if isOwner, _ := functions.IsSalonOwner(userId, quote.SalonId); !isOwner {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Access denied: you are not the owner of this salon")
	}

	quote.QuoteId = uuid.New().String()
	now := time.Now()
	quote.CreatedAt = now
	quote.UpdatedAt = now
	quote.Deleted = false

	if err := dao.CreateQuote(quote); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create quote")
	}
	return utils.SendSuccessResponse(c)
}
