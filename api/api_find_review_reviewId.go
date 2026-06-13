package api

import (
	"net/http"
	"strings"

	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func FindReviewApi(c *fiber.Ctx) error {
	reviewId := c.Query("ReviewId")
	if reviewId == "" {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "ReviewId is required")
	}

	review, err := dao.FindReviewByReviewId(reviewId)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusNotFound, "review not found")
	}

	if !review.IsVisible {
		authorized := false
		authorization := c.Get("Authorization")
		if authorization != "" {
			parts := strings.Fields(authorization)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				claims, err := functions.ValidateJWT(parts[1])
				if err == nil {
					isAdminOrMod := claims.Role == "admin" || claims.Role == "moderator"
					if review.UserId == claims.UserId || isAdminOrMod {
						authorized = true
					}
				}
			}
		}
		if !authorized {
			return utils.SendErrorResponse(c, http.StatusNotFound, "Review not found")
		}
	}

	return utils.SendDataResponse(c, review)
}
