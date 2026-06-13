package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindScheduleBlocksByBarberDate returns all blocks for a barber that overlap a
// calendar date. The date is treated as midnight-to-midnight UTC so it works
// consistently across time zones for the legacy slot-generation flow.
// Pass empty date to return all blocks for the barber.
func FindScheduleBlocksByBarberDate(barberId, date string) ([]dto.ScheduleBlock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := dbConfig.DATABASE.Collection(dbConfig.SCHEDULE_BLOCKS_COLLECTION)

	filter := bson.M{"BarberId": barberId}
	if date != "" {
		dayStart, err := time.Parse("2006-01-02", date)
		if err == nil {
			dayEnd := dayStart.AddDate(0, 0, 1)
			// overlap: block.StartTime < dayEnd AND block.EndTime > dayStart
			filter["StartTime"] = bson.M{"$lt": dayEnd}
			filter["EndTime"] = bson.M{"$gt": dayStart}
		}
	}

	cursor, err := collection.Find(ctx, filter, options.Find().SetSort(bson.M{"StartTime": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var blocks []dto.ScheduleBlock
	if err := cursor.All(ctx, &blocks); err != nil {
		return nil, err
	}
	return blocks, nil
}
