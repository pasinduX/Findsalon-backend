package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func CreateWorkingHours(hours dto.WorkingHours) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.HOURS_COLLECTION)
	_, err := collection.InsertOne(ctx, hours)
	return err
}
