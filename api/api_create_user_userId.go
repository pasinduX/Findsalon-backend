package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateUserApi(c *fiber.Ctx) error {
	var user dto.User
	if err := c.BodyParser(&user); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if !functions.UniqueCheck(dbConfig.USERS_COLLECTION, "Email", user.Email) {
		return utils.SendErrorResponse(c, fiber.StatusConflict, "Email already exists")
	}

	user.UserId = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = true
	user.Deleted = false

	if user.Role == "" {
		user.Role = "user"
	}

	if user.Provider == "" {
		user.Provider = "local"
	}

	if user.PasswordHash != "" {
		hashedPassword, err := functions.HashPassword(user.PasswordHash)
		if err != nil {
			return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to hash password")
		}
		user.PasswordHash = hashedPassword
	}

	if err := dao.CreateUser(user); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user")
	}

	return utils.SendSuccessResponse(c)
}
