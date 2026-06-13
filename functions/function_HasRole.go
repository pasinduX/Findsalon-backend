package functions

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func HasRole(userId, role string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := dbConfig.DATABASE.Collection(dbConfig.USERROLES_COLLECTION)
	count, err := collection.CountDocuments(ctx, bson.M{"UserId": userId, "Role": role, "Deleted": false})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func HasRoleInSalon(userId, role, salonId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := dbConfig.DATABASE.Collection(dbConfig.USERROLES_COLLECTION)
	count, err := collection.CountDocuments(ctx, bson.M{"UserId": userId, "Role": role, "SalonId": salonId, "Deleted": false})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetUserRoles(userId string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := dbConfig.DATABASE.Collection(dbConfig.USERROLES_COLLECTION)
	cursor, err := collection.Find(ctx, bson.M{"UserId": userId, "Deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var roles []string
	for cursor.Next(ctx) {
		var r dto.UserRole
		if err := cursor.Decode(&r); err != nil {
			continue
		}
		roles = append(roles, r.Role)
	}
	return roles, nil
}

func AdminOnlyMiddleware(c *fiber.Ctx) error {
	role, _ := c.Locals("role").(string)
	if role != dto.RoleAdmin {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "admin privileges required")
	}
	return c.Next()
}

func AdminOrModeratorMiddleware(c *fiber.Ctx) error {
	role, _ := c.Locals("role").(string)
	if role != dto.RoleAdmin && role != dto.Rolemoderator {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "admin or moderator privileges required")
	}
	return c.Next()
}
