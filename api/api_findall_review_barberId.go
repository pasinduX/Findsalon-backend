package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
)

func FindAllReviewByBarberApi(c *fiber.Ctx) error {
	barberId := c.Query("BarberId")
	if barberId == "" {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "BarberId is required")
	}

	pagination := functions.GetPaginationParams(c)
	skip, limit := functions.BuildSkipLimit(pagination)

	reviews, total, err := dao.FindAllReviewsByBarberId(barberId, skip, limit)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, "failed to retrieve reviews")
	}

	response := functions.BuildPaginatedResponse(reviews, total, pagination)
	return utils.SendDataResponse(c, response)
}
