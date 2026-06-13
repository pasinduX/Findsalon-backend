package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func FindAllReviewsBySalonId(salonId string, skip int64, limit int64) ([]dto.Review, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"SalonId": salonId, "Deleted": false, "IsVisible": true}
	collection := dbConfig.DATABASE.Collection(dbConfig.REVIEWS_COLLECTION)

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := collection.Find(ctx, filter, options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{"CreatedAt", -1}}))
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var reviews []dto.Review
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, 0, err
	}

	return reviews, int(total), nil
}
