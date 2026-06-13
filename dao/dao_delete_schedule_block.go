package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"

	"go.mongodb.org/mongo-driver/bson"
)

func DeleteScheduleBlock(blockId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := dbConfig.DATABASE.Collection(dbConfig.SCHEDULE_BLOCKS_COLLECTION)
	_, err := collection.DeleteOne(ctx, bson.M{"BlockId": blockId})
	return err
}

// DeleteNonBookedSlotsForDateBarber removes all non-booked slots for a barber on a date
// so fresh slots can be regenerated.
func DeleteNonBookedSlotsForDateBarber(barberId, date string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := dbConfig.DATABASE.Collection(dbConfig.TIMESLOTS_COLLECTION)
	_, err := collection.DeleteMany(ctx, bson.M{
		"BarberId": barberId,
		"Date":     date,
		"IsBooked": false,
	})
	return err
}
