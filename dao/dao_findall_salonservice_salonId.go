package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func FindAllServicesBySalonId(salonId string) ([]dto.SalonService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.SALONSERVICES_COLLECTION)
	opts := options.Find().SetSort(bson.M{"SortOrder": 1})
	cursor, err := collection.Find(ctx, bson.M{"SalonId": salonId, "Deleted": false}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var services []dto.SalonService
	if err = cursor.All(ctx, &services); err != nil {
		return nil, err
	}
	return services, nil
}
