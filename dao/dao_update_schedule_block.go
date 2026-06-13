package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateScheduleBlock(blockId string, update dto.ScheduleBlock) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := dbConfig.DATABASE.Collection(dbConfig.SCHEDULE_BLOCKS_COLLECTION)
	update.UpdatedAt = time.Now()
	_, err := collection.UpdateOne(ctx,
		bson.M{"BlockId": blockId},
		bson.M{"$set": update},
	)
	return err
}
