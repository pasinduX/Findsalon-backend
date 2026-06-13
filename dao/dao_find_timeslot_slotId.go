package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"go.mongodb.org/mongo-driver/bson"
)

func FindTimeSlotBySlotId(slotId string) (dto.TimeSlot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.TIMESLOTS_COLLECTION)

	var slot dto.TimeSlot
	err := collection.FindOne(ctx, bson.M{"SlotId": slotId, "Deleted": false}).Decode(&slot)
	return slot, err
}
