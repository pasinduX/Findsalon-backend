package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"

	"go.mongodb.org/mongo-driver/bson"
)

func MarkTimeSlotBooked(slotId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.TIMESLOTS_COLLECTION)
	_, err := collection.UpdateOne(ctx,
		bson.M{"SlotId": slotId},
		bson.M{"$set": bson.M{
			"IsBooked":  true,
			"UpdatedAt": time.Now(),
		}},
	)
	return err
}

func MarkTimeSlotAvailable(slotId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.TIMESLOTS_COLLECTION)
	_, err := collection.UpdateOne(ctx,
		bson.M{"SlotId": slotId},
		bson.M{"$set": bson.M{
			"IsBooked":  false,
			"UpdatedAt": time.Now(),
		}},
	)
	return err
}
