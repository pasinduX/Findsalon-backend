package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindAllTimeSlotsByBarberId(barberId string, filter bson.M) ([]dto.TimeSlot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.TIMESLOTS_COLLECTION)

	query := bson.M{"BarberId": barberId, "Deleted": false}
	for k, v := range filter {
		query[k] = v
	}

	opts := options.Find().SetSort(bson.D{
		{Key: "Date", Value: 1},
		{Key: "StartTime", Value: 1},
	})

	cursor, err := collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var slots []dto.TimeSlot
	if err := cursor.All(ctx, &slots); err != nil {
		return nil, err
	}
	return slots, nil
}
