package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindAllBookingsByUserId(userId string, filter bson.M) ([]dto.Booking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.BOOKINGS_COLLECTION)

	query := bson.M{"UserId": userId, "Deleted": false}
	for k, v := range filter {
		query[k] = v
	}

	opts := options.Find().SetSort(bson.D{{Key: "CreatedAt", Value: -1}})

	cursor, err := collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []dto.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}
