package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"findsalon-backend/dao"
	"findsalon-backend/utils"
)

func FindBarberRatingSummaryApi(c *fiber.Ctx) error {
	barberId := c.Query("BarberId")
	if barberId == "" {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "BarberId is required")
	}

	summary, err := dao.AggregateRatingSummary("BarberId", barberId)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, "failed to build rating summary")
	}

	return utils.SendDataResponse(c, summary)
}
