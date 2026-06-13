package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"findsalon-backend/dao"
	"findsalon-backend/utils"
)

func DownloadReviewApi(c *fiber.Ctx) error {
	filter := bson.M{}
	if salonId := c.Query("SalonId"); salonId != "" {
		filter["SalonId"] = salonId
	}
	if barberId := c.Query("BarberId"); barberId != "" {
		filter["BarberId"] = barberId
	}

	if fromDate := c.Query("FromDate"); fromDate != "" {
		parsedFrom, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			return utils.SendErrorResponse(c, http.StatusBadRequest, "invalid FromDate format")
		}
		if filter["CreatedAt"] == nil {
			filter["CreatedAt"] = bson.M{}
		}
		filter["CreatedAt"].(bson.M)["$gte"] = parsedFrom
	}

	if toDate := c.Query("ToDate"); toDate != "" {
		parsedTo, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			return utils.SendErrorResponse(c, http.StatusBadRequest, "invalid ToDate format")
		}
		if filter["CreatedAt"] == nil {
			filter["CreatedAt"] = bson.M{}
		}
		filter["CreatedAt"].(bson.M)["$lte"] = parsedTo
	}

	reviews, _, err := dao.FindAllReviews(filter, 0, 100000)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, "failed to retrieve reviews")
	}

	content, err := json.Marshal(reviews)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, "failed to marshal reviews")
	}

	filename := "reviews_export_" + time.Now().Format("20060102") + ".json"
	c.Set("Content-Type", "application/json")
	c.Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	return c.Status(http.StatusOK).Send(content)
}
