package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"findsalon-backend/dbConfig"
)

func DeleteSalonService(serviceId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.SALONSERVICES_COLLECTION)
	update := bson.M{
		"$set": bson.M{
			"Deleted":   true,
			"UpdatedAt": time.Now(),
		},
	}
	_, err := collection.UpdateOne(ctx, bson.M{"ServiceId": serviceId}, update)
	return err
}
