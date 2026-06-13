package functions

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"go.mongodb.org/mongo-driver/bson"
)

// FindBookingForReview fetches a booking and validates it belongs to userId and is completed.
// Replaces the HTTP call to booking-service in the old review-ms.
func FindBookingForReview(bookingId, userId string) (*dto.Booking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.BOOKINGS_COLLECTION)
	var booking dto.Booking
	if err := col.FindOne(ctx, bson.M{"BookingId": bookingId, "Deleted": false}).Decode(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}

// FetchSalonForReview returns basic salon info for a review (validates the salon exists).
func FetchSalonForReview(salonId string) (*dto.Salon, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.SALONS_COLLECTION)
	var salon dto.Salon
	if err := col.FindOne(ctx, bson.M{"SalonId": salonId, "Deleted": false}).Decode(&salon); err != nil {
		return nil, err
	}
	return &salon, nil
}

// FetchBarberForReview validates the barber exists.
func FetchBarberForReview(barberId string) (*dto.Barber, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.BARBERS_COLLECTION)
	var barber dto.Barber
	if err := col.FindOne(ctx, bson.M{"BarberId": barberId, "Deleted": false}).Decode(&barber); err != nil {
		return nil, err
	}
	return &barber, nil
}
