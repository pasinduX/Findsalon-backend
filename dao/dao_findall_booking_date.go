package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindAllBookingsByDate(date string, salonId string) ([]dto.Booking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slotCollection := dbConfig.DATABASE.Collection(dbConfig.TIMESLOTS_COLLECTION)
	slotCursor, err := slotCollection.Find(ctx, bson.M{"Date": date, "SalonId": salonId, "Deleted": false})
	if err != nil {
		return nil, err
	}
	defer slotCursor.Close(ctx)

	var slots []dto.TimeSlot
	if err := slotCursor.All(ctx, &slots); err != nil {
		return nil, err
	}

	slotIds := make([]string, 0, len(slots))
	for _, s := range slots {
		slotIds = append(slotIds, s.SlotId)
	}

	if len(slotIds) == 0 {
		return []dto.Booking{}, nil
	}

	bookingCollection := dbConfig.DATABASE.Collection(dbConfig.BOOKINGS_COLLECTION)

	bookingCtx, bookingCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer bookingCancel()

	opts := options.Find().SetSort(bson.D{{Key: "CreatedAt", Value: -1}})
	bookingCursor, err := bookingCollection.Find(bookingCtx, bson.M{
		"SlotId":  bson.M{"$in": slotIds},
		"SalonId": salonId,
		"Deleted": false,
	}, opts)
	if err != nil {
		return nil, err
	}
	defer bookingCursor.Close(bookingCtx)

	var bookings []dto.Booking
	if err := bookingCursor.All(bookingCtx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}
