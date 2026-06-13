package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func LogoutApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)

	if err := dao.UpdateUser(userId, bson.M{"RefreshToken": ""}); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to logout user")
	}

	return utils.SendSuccessResponse(c)
}
