package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
)

func FindAllReviewBySalonApi(c *fiber.Ctx) error {
	salonId := c.Query("SalonId")
	if salonId == "" {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "SalonId is required")
	}

	pagination := functions.GetPaginationParams(c)
	skip, limit := functions.BuildSkipLimit(pagination)

	reviews, total, err := dao.FindAllReviewsBySalonId(salonId, skip, limit)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, "failed to retrieve reviews")
	}

	response := functions.BuildPaginatedResponse(reviews, total, pagination)
	return utils.SendDataResponse(c, response)
}
