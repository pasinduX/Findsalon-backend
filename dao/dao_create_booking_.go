package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func CreateBooking(booking dto.Booking) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.BOOKINGS_COLLECTION)
	_, err := collection.InsertOne(ctx, booking)
	return err
}
