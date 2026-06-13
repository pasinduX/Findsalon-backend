package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"findsalon-backend/dbConfig"
)

func UpdateBarber(barberId string, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update["UpdatedAt"] = time.Now()
	collection := dbConfig.DATABASE.Collection(dbConfig.BARBERS_COLLECTION)
	_, err := collection.UpdateOne(ctx, bson.M{"BarberId": barberId}, bson.M{"$set": update})
	return err
}
