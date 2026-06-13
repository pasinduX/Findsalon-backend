package dao

import (
	"context"
	"errors"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"go.mongodb.org/mongo-driver/bson"
)

func DeleteTimeSlot(slotId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.TIMESLOTS_COLLECTION)

	var slot dto.TimeSlot
	err := collection.FindOne(ctx, bson.M{"SlotId": slotId, "Deleted": false}).Decode(&slot)
	if err != nil {
		return err
	}

	if slot.IsBooked {
		return errors.New("cannot delete a booked time slot")
	}

	_, err = collection.UpdateOne(ctx,
		bson.M{"SlotId": slotId, "Deleted": false},
		bson.M{"$set": bson.M{
			"Deleted":   true,
			"UpdatedAt": time.Now(),
		}},
	)
	return err
}
