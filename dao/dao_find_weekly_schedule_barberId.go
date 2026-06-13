package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"go.mongodb.org/mongo-driver/bson"
)

func FindWeeklyScheduleByBarber(barberId string) (*dto.WeeklySchedule, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := dbConfig.DATABASE.Collection(dbConfig.WEEKLY_SCHEDULES_COLLECTION)
	var schedule dto.WeeklySchedule
	err := collection.FindOne(ctx, bson.M{"BarberId": barberId}).Decode(&schedule)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}
