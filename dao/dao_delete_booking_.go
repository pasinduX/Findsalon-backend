package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"

	"go.mongodb.org/mongo-driver/bson"
)

func DeleteBooking(bookingId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.BOOKINGS_COLLECTION)
	_, err := collection.UpdateOne(ctx,
		bson.M{"BookingId": bookingId, "Deleted": false},
		bson.M{"$set": bson.M{
			"Deleted":   true,
			"UpdatedAt": time.Now(),
		}},
	)
	return err
}
