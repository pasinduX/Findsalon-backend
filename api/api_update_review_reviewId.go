package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"findsalon-backend/dao"
	"findsalon-backend/utils"
)

func UpdateReviewApi(c *fiber.Ctx) error {
	reviewId := c.Query("ReviewId")
	if reviewId == "" {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "ReviewId is required")
	}

	review, err := dao.FindReviewByReviewId(reviewId)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusNotFound, "review not found")
	}

	userId, _ := c.Locals("userId").(string)
	role, _ := c.Locals("role").(string)
	if review.UserId != userId && role != "admin" && role != "moderator" {
		return utils.SendErrorResponse(c, http.StatusForbidden, "not authorized to update this review")
	}

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	allowed := bson.M{}
	if value, found := body["Rating"]; found {
		rating, ok := value.(float64)
		if !ok || int(rating) < 1 || int(rating) > 5 {
			return utils.SendErrorResponse(c, http.StatusBadRequest, "Rating must be between 1 and 5")
		}
		allowed["Rating"] = int(rating)
	}
	if value, found := body["Comment"]; found {
		if comment, ok := value.(string); ok {
			allowed["Comment"] = comment
		}
	}
	if value, found := body["IsVisible"]; found {
		if visible, ok := value.(bool); ok {
			allowed["IsVisible"] = visible
		}
	}

	if len(allowed) == 0 {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "no updatable fields provided")
	}

	if err := dao.UpdateReview(reviewId, allowed); err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, "failed to update review")
	}

	return utils.SendSuccessResponse(c)
}
