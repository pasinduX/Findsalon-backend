package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"findsalon-backend/dao"
	"findsalon-backend/utils"
)

func FindSalonRatingSummaryApi(c *fiber.Ctx) error {
	salonId := c.Query("SalonId")
	if salonId == "" {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "SalonId is required")
	}

	summary, err := dao.AggregateRatingSummary("SalonId", salonId)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, "failed to build rating summary")
	}

	return utils.SendDataResponse(c, summary)
}
