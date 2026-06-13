package api

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func isValidRole(role string) bool {
	switch role {
	case dto.RoleAdmin, dto.Rolemoderator, dto.RoleUser, dto.RoleSalonOwner, dto.RoleBarber:
		return true
	}
	return false
}

// isSystemRole returns true for roles that are not salon-scoped.
func isSystemRole(role string) bool {
	return role == dto.RoleAdmin || role == dto.Rolemoderator || role == dto.RoleUser
}

func CreateUserRoleApi(c *fiber.Ctx) error {

	var body dto.UserRole
	if err := c.BodyParser(&body); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}
	v := validator.New()
	if err := v.Struct(body); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "validation failed")
	}
	if !isValidRole(body.Role) {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "invalid role value")
	}

	if !isSystemRole(body.Role) && body.SalonId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SalonId is required for salon-scoped roles")
	}

	existingRoles, err := dao.FindUserRolesByUserId(body.UserId)
	if err == nil {
		for _, existing := range existingRoles {
			if existing.Role == body.Role && !existing.Deleted {
				if isSystemRole(body.Role) || existing.SalonId == body.SalonId {
					return utils.SendErrorResponse(c, fiber.StatusConflict, "User already has this role")
				}
			}
		}
	}

	body.RoleId = uuid.New().String()
	body.CreatedAt = time.Now()
	body.UpdatedAt = time.Now()
	body.Deleted = false

	if err := dao.CreateUserRole(body); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to create user role")
	}

	// Keep the User.Role field in sync for JWT claims.
	_ = dao.UpdateUser(body.UserId, bson.M{"Role": body.Role})

	return utils.SendSuccessResponse(c)
}
