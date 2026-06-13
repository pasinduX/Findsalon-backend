package api

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateUserRoleApi(c *fiber.Ctx) error {

	userId := c.Query("UserId")
	if userId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "UserId is required")
	}

	var payload map[string]string
	if err := c.BodyParser(&payload); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}
	newRole, ok := payload["Role"]
	if !ok || newRole == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Role is required")
	}
	if !isValidRole(newRole) {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "invalid role value")
	}

	currentRoles, err := dao.FindUserRolesByUserId(userId)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to fetch current roles")
	}
	for _, role := range currentRoles {
		if err := dao.DeleteUserRole(role.RoleId); err != nil {
			log.Printf("failed to soft delete role %s: %v", role.RoleId, err)
		}
	}

	newRecord := dto.UserRole{
		RoleId:    uuid.New().String(),
		UserId:    userId,
		Role:      newRole,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Deleted:   false,
	}
	if err := dao.CreateUserRole(newRecord); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to create updated user role")
	}

	// Keep User.Role in sync for JWT claims.
	_ = dao.UpdateUser(userId, bson.M{"Role": newRole})

	return utils.SendSuccessResponse(c)
}
