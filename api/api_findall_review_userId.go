package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
)

func FindAllReviewByUserApi(c *fiber.Ctx) error {
	userId := c.Query("UserId")
	currentUserId, _ := c.Locals("userId").(string)
	role, _ := c.Locals("role").(string)

	if userId == "" {
		userId = currentUserId
	}

	if userId != currentUserId && role != "admin" && role != "moderator" {
		return utils.SendErrorResponse(c, http.StatusForbidden, "not authorized to view this user's reviews")
	}

	pagination := functions.GetPaginationParams(c)
	skip, limit := functions.BuildSkipLimit(pagination)

	reviews, total, err := dao.FindAllReviewsByUserId(userId, skip, limit)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, "failed to retrieve reviews")
	}

	response := functions.BuildPaginatedResponse(reviews, total, pagination)
	return utils.SendDataResponse(c, response)
}
