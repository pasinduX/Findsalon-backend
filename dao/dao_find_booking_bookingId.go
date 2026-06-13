package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"go.mongodb.org/mongo-driver/bson"
)

func FindBookingByBookingId(bookingId string) (dto.Booking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.BOOKINGS_COLLECTION)

	var booking dto.Booking
	err := collection.FindOne(ctx, bson.M{"BookingId": bookingId, "Deleted": false}).Decode(&booking)
	return booking, err
}
