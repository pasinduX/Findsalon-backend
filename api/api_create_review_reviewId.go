package api

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
)

func CreateReviewApi(c *fiber.Ctx) error {
	var review dto.Review
	if err := c.BodyParser(&review); err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	v := validator.New()
	if err := v.Struct(&review); err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	userId, ok := c.Locals("userId").(string)
	if !ok || userId == "" {
		return utils.SendErrorResponse(c, http.StatusUnauthorized, "invalid user context")
	}
	review.UserId = userId

	if _, err := functions.FetchSalonForReview(review.SalonId); err != nil {
		return utils.SendErrorResponse(c, http.StatusNotFound, "Salon not found")
	}

	if review.BarberId != "" {
		barber, err := functions.FetchBarberForReview(review.BarberId)
		if err != nil || barber.SalonId != review.SalonId {
			return utils.SendErrorResponse(c, http.StatusBadRequest, "Barber does not belong to this salon")
		}
	}

	if review.BookingId != "" {
		booking, err := functions.FindBookingForReview(review.BookingId, userId)
		if err != nil {
			return utils.SendErrorResponse(c, http.StatusInternalServerError, "unable to validate booking")
		}
		if booking.Status != dto.BookingStatusCompleted {
			return utils.SendErrorResponse(c, http.StatusBadRequest, "Can only review a completed booking")
		}
		if booking.UserId != userId {
			return utils.SendErrorResponse(c, http.StatusForbidden, "You can only review your own booking")
		}
		if _, err := dao.FindReviewByBookingId(review.BookingId); err == nil {
			return utils.SendErrorResponse(c, http.StatusConflict, "A review already exists for this booking")
		}
	}

	// Populate UserName from Users collection directly.
	review.UserName = "Anonymous"
	if user, err := dao.FindUserByUserId(userId); err == nil {
		review.UserName = user.FullName
	}

	review.ReviewId = uuid.New().String()
	review.IsVisible = true
	review.CreatedAt = time.Now()
	review.UpdatedAt = review.CreatedAt
	review.Deleted = false

	if err := dao.CreateReview(review); err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, "failed to create review")
	}

	return utils.SendSuccessResponse(c)
}
