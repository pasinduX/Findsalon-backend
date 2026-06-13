package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"findsalon-backend/dbConfig"
)

func DeleteBarber(barberId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.BARBERS_COLLECTION)
	update := bson.M{
		"$set": bson.M{
			"Deleted":   true,
			"IsActive":  false,
			"UpdatedAt": time.Now(),
		},
	}
	_, err := collection.UpdateOne(ctx, bson.M{"BarberId": barberId}, update)
	return err
}
