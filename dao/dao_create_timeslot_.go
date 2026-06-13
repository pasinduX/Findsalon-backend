package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func CreateTimeSlot(slot dto.TimeSlot) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.TIMESLOTS_COLLECTION)
	_, err := collection.InsertOne(ctx, slot)
	return err
}
