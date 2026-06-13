package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"findsalon-backend/dao"
	"findsalon-backend/utils"
)

func DeleteReviewApi(c *fiber.Ctx) error {
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
	isAdminOrMod := role == "admin" || role == "moderator"
	if review.UserId != userId && !isAdminOrMod {
		return utils.SendErrorResponse(c, http.StatusForbidden, "not authorized to delete this review")
	}

	if err := dao.DeleteReview(reviewId); err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, "failed to delete review")
	}

	return utils.SendSuccessResponse(c)
}
