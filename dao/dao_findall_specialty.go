package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"go.mongodb.org/mongo-driver/bson"
)

func FindAllSpecialties() ([]dto.Specialty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := dbConfig.DATABASE.Collection(dbConfig.SPECIALTIES_COLLECTION)
	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var specialties []dto.Specialty
	if err = cursor.All(ctx, &specialties); err != nil {
		return nil, err
	}
	return specialties, nil
}
