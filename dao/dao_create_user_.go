package dao

import (
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
	"context"
	"time"
)

func CreateUser(user dto.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.USERS_COLLECTION)
	_, err := collection.InsertOne(ctx, user)
	return err
}
