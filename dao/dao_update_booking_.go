package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateBooking(bookingId string, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if setFields, ok := update["$set"]; ok {
		if setMap, ok := setFields.(bson.M); ok {
			setMap["UpdatedAt"] = time.Now()
		}
	} else {
		update["$set"] = bson.M{"UpdatedAt": time.Now()}
	}

	collection := dbConfig.DATABASE.Collection(dbConfig.BOOKINGS_COLLECTION)
	_, err := collection.UpdateOne(ctx,
		bson.M{"BookingId": bookingId, "Deleted": false},
		update,
	)
	return err
}
