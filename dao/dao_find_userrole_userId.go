package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func FindUserRolesByUserId(userId string) ([]dto.UserRole, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.USERROLES_COLLECTION)
	cursor, err := collection.Find(ctx, bson.M{"UserId": userId, "Deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	roles := []dto.UserRole{}
	for cursor.Next(ctx) {
		var role dto.UserRole
		if err := cursor.Decode(&role); err != nil {
			continue
		}
		roles = append(roles, role)
	}
	return roles, nil
}
