package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func FindReviewByBookingId(bookingId string) (dto.Review, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var review dto.Review
	collection := dbConfig.DATABASE.Collection(dbConfig.REVIEWS_COLLECTION)
	if err := collection.FindOne(ctx, bson.M{"BookingId": bookingId, "Deleted": false}).Decode(&review); err != nil {
		return dto.Review{}, err
	}
	return review, nil
}
