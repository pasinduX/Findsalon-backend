package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"findsalon-backend/dbConfig"
)

func UpdateWorkingHours(hoursId string, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update["UpdatedAt"] = time.Now()
	collection := dbConfig.DATABASE.Collection(dbConfig.HOURS_COLLECTION)
	_, err := collection.UpdateOne(ctx, bson.M{"HoursId": hoursId}, bson.M{"$set": update})
	return err
}
