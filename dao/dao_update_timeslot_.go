package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateTimeSlot(slotId string, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if setFields, ok := update["$set"]; ok {
		if setMap, ok := setFields.(bson.M); ok {
			setMap["UpdatedAt"] = time.Now()
		}
	} else {
		update["$set"] = bson.M{"UpdatedAt": time.Now()}
	}

	collection := dbConfig.DATABASE.Collection(dbConfig.TIMESLOTS_COLLECTION)
	_, err := collection.UpdateOne(ctx,
		bson.M{"SlotId": slotId, "Deleted": false},
		update,
	)
	return err
}
