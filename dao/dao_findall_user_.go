package dao

import (
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func FindAllUsers() ([]dto.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.USERS_COLLECTION)
	cursor, err := collection.Find(ctx, bson.M{"Deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []dto.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
