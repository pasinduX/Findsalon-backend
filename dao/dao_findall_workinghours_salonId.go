package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func FindAllWorkingHoursBySalonId(salonId string) ([]dto.WorkingHours, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.HOURS_COLLECTION)
	opts := options.Find().SetSort(bson.M{"SortOrder": 1})
	cursor, err := collection.Find(ctx, bson.M{"SalonId": salonId, "Deleted": false}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var hours []dto.WorkingHours
	if err = cursor.All(ctx, &hours); err != nil {
		return nil, err
	}
	return hours, nil
}
