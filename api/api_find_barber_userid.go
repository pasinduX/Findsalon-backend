package api

import (
	"errors"

	"findsalon-backend/dao"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindBarberByUserApi(c *fiber.Ctx) error {

	userId, _ := c.Locals("userId").(string)
	email, _ := c.Locals("email").(string)

	barber, err := dao.FindBarberByUserId(userId)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve barber profile")
		}

		// Not linked yet — try to find by email and auto-link
		if email != "" {
			barber, err = dao.FindBarberByEmail(email)
			if err == nil {
				// Auto-link: persist the userId onto the barber record
				_ = dao.UpdateBarberUserId(barber.BarberId, userId)
				barber.UserId = userId
				return utils.SendDataResponse(c, barber)
			}
		}

		// No barber profile — return 200 with null so clients treat it as
		// "not a barber" rather than a network error (avoids noisy 404s).
		return utils.SendDataResponse(c, nil)
	}

	return utils.SendDataResponse(c, barber)
}
