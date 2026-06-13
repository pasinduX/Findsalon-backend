package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"findsalon-backend/dbConfig"
)

func UpdateSalon(salonId string, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update["UpdatedAt"] = time.Now()
	collection := dbConfig.DATABASE.Collection(dbConfig.SALONS_COLLECTION)
	_, err := collection.UpdateOne(ctx, bson.M{"SalonId": salonId}, bson.M{"$set": update})
	return err
}
