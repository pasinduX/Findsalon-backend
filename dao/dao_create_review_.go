package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

// db.Reviews.createIndex({ BookingId: 1 }, { unique: true, sparse: true })
func CreateReview(review dto.Review) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.REVIEWS_COLLECTION)
	_, err := collection.InsertOne(ctx, review)
	return err
}
