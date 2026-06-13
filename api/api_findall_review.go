package api

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
)

func FindAllReviewApi(c *fiber.Ctx) error {
	pagination := functions.GetPaginationParams(c)
	skip, limit := functions.BuildSkipLimit(pagination)

	filter := bson.M{}
	if salonId := c.Query("SalonId"); salonId != "" {
		filter["SalonId"] = salonId
	}
	if barberId := c.Query("BarberId"); barberId != "" {
		filter["BarberId"] = barberId
	}
	if userId := c.Query("UserId"); userId != "" {
		filter["UserId"] = userId
	}
	if rating := c.Query("Rating"); rating != "" {
		if value, err := strconv.Atoi(rating); err == nil {
			filter["Rating"] = value
		}
	}
	if isVisible := c.Query("IsVisible"); isVisible != "" {
		if value, err := strconv.ParseBool(isVisible); err == nil {
			filter["IsVisible"] = value
		}
	}

	reviews, total, err := dao.FindAllReviews(filter, skip, limit)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, "failed to retrieve reviews")
	}

	response := functions.BuildPaginatedResponse(reviews, total, pagination)
	return utils.SendDataResponse(c, response)
}
