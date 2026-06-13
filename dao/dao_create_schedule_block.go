package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func CreateScheduleBlock(block dto.ScheduleBlock) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := dbConfig.DATABASE.Collection(dbConfig.SCHEDULE_BLOCKS_COLLECTION)
	_, err := collection.InsertOne(ctx, block)
	return err
}
