package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func FindAllSalons(filters bson.M) ([]dto.Salon, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"Deleted": false}
	for k, v := range filters {
		filter[k] = v
	}

	collection := dbConfig.DATABASE.Collection(dbConfig.SALONS_COLLECTION)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var salons []dto.Salon
	if err = cursor.All(ctx, &salons); err != nil {
		return nil, err
	}
	return salons, nil
}
