package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func CreateWeeklySchedule(schedule dto.WeeklySchedule) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := dbConfig.DATABASE.Collection(dbConfig.WEEKLY_SCHEDULES_COLLECTION)
	_, err := collection.InsertOne(ctx, schedule)
	return err
}
