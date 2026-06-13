package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"go.mongodb.org/mongo-driver/bson"
)

func FindScheduleBlockById(blockId string) (*dto.ScheduleBlock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := dbConfig.DATABASE.Collection(dbConfig.SCHEDULE_BLOCKS_COLLECTION)
	var block dto.ScheduleBlock
	err := collection.FindOne(ctx, bson.M{"BlockId": blockId}).Decode(&block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}
