package functions

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"

	"go.mongodb.org/mongo-driver/bson"
)

func IsSlotAvailable(slotId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.TIMESLOTS_COLLECTION)
	count, err := col.CountDocuments(ctx, bson.M{"SlotId": slotId, "IsBooked": false, "Deleted": false})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func ValidateSlotOwnership(slotId, barberId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.TIMESLOTS_COLLECTION)
	count, err := col.CountDocuments(ctx, bson.M{"SlotId": slotId, "BarberId": barberId, "Deleted": false})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
