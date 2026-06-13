package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"findsalon-backend/dbConfig"
)

func UpdateSalonService(serviceId string, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update["UpdatedAt"] = time.Now()
	collection := dbConfig.DATABASE.Collection(dbConfig.SALONSERVICES_COLLECTION)
	_, err := collection.UpdateOne(ctx, bson.M{"ServiceId": serviceId}, bson.M{"$set": update})
	return err
}
