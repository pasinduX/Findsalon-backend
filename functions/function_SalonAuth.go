package functions

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"go.mongodb.org/mongo-driver/bson"
)

// IsSalonOwner returns true when userId matches the OwnerId of the given salon.
// Returns true unconditionally when userId is empty (auth disabled).
func IsSalonOwner(userId, salonId string) (bool, error) {
	if userId == "" {
		return true, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.SALONS_COLLECTION)
	var salon dto.Salon
	err := col.FindOne(ctx, bson.M{"SalonId": salonId, "Deleted": false}).Decode(&salon)
	if err != nil {
		return false, nil
	}
	return salon.OwnerId == userId, nil
}

// IsBarberOfSalon returns true when a Barber record for this userId+salonId exists.
func IsBarberOfSalon(userId, salonId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.BARBERS_COLLECTION)
	count, err := col.CountDocuments(ctx, bson.M{"UserId": userId, "SalonId": salonId, "Deleted": false})
	if err != nil {
		return false, nil
	}
	return count > 0, nil
}

// IsBarberOfSlot returns true when a time slot's owning barber has UserId == userId.
func IsBarberOfSlot(userId, slotId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	slotCol := dbConfig.DATABASE.Collection(dbConfig.TIMESLOTS_COLLECTION)
	var slot dto.TimeSlot
	if err := slotCol.FindOne(ctx, bson.M{"SlotId": slotId, "Deleted": false}).Decode(&slot); err != nil {
		return false, nil
	}
	return IsBarberOfSalon(userId, slot.SalonId)
}

// IsSalonOwnerOrBarber returns true if userId is either the salon owner or any barber in it.
// Returns true unconditionally when userId is empty (auth disabled).
func IsSalonOwnerOrBarber(userId, salonId string) (bool, error) {
	if userId == "" {
		return true, nil
	}
	isOwner, _ := IsSalonOwner(userId, salonId)
	if isOwner {
		return true, nil
	}
	return IsBarberOfSalon(userId, salonId)
}

// IsBarberUser returns true when any barber record in the DB has UserId == userId.
func IsBarberUser(userId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.BARBERS_COLLECTION)
	count, err := col.CountDocuments(ctx, bson.M{"UserId": userId, "Deleted": false})
	if err != nil {
		return false, nil
	}
	return count > 0, nil
}
