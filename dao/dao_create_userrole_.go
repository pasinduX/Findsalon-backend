package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func CreateUserRole(role dto.UserRole) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.USERROLES_COLLECTION)
	_, err := collection.InsertOne(ctx, role)
	return err
}
