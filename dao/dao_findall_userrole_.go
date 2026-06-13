package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func FindAllUserRoles(skip int64, limit int64) ([]dto.UserRole, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.USERROLES_COLLECTION)
	filter := bson.M{"Deleted": false}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	options := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.M{"CreatedAt": -1})
	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, 0, err
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
	return roles, int(count), nil
}
