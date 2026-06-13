package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func BulkCreateTimeSlots(slots []dto.TimeSlot) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.TIMESLOTS_COLLECTION)

	docs := make([]interface{}, len(slots))
	for i, slot := range slots {
		docs[i] = slot
	}

	_, err := collection.InsertMany(ctx, docs)
	return err
}
