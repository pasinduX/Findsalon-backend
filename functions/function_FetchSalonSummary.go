package functions

import (
	"context"
	"log"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FetchUserSalons returns all salons owned by a user.
// Replaces the HTTP call to salon-service in the old user-ms.
func FetchUserSalons(ownerId string) ([]dto.SalonSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.SALONS_COLLECTION)
	cursor, err := col.Find(ctx, bson.M{"OwnerId": ownerId, "Deleted": false})
	if err != nil {
		return []dto.SalonSummary{}, nil
	}
	defer cursor.Close(ctx)
	var summaries []dto.SalonSummary
	for cursor.Next(ctx) {
		var salon dto.Salon
		if err := cursor.Decode(&salon); err != nil {
			continue
		}
		summaries = append(summaries, dto.SalonSummary{
			SalonId:  salon.SalonId,
			OwnerId:  salon.OwnerId,
			Name:     salon.Name,
			IsActive: salon.IsActive,
		})
	}
	if summaries == nil {
		summaries = []dto.SalonSummary{}
	}
	return summaries, nil
}

// FetchSalonSummaryById returns a single salon summary.
func FetchSalonSummaryById(salonId string) (dto.SalonSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.SALONS_COLLECTION)
	var salon dto.Salon
	if err := col.FindOne(ctx, bson.M{"SalonId": salonId, "Deleted": false}).Decode(&salon); err != nil {
		return dto.SalonSummary{}, err
	}
	return dto.SalonSummary{SalonId: salon.SalonId, OwnerId: salon.OwnerId, Name: salon.Name, IsActive: salon.IsActive}, nil
}

// FetchBarberMemberships returns all salon memberships for a barber user.
func FetchBarberMemberships(userId string) ([]dto.SalonMembership, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.BARBERS_COLLECTION)
	cursor, err := col.Find(ctx, bson.M{"UserId": userId, "Deleted": false})
	if err != nil {
		return []dto.SalonMembership{}, nil
	}
	defer cursor.Close(ctx)
	var memberships []dto.SalonMembership
	for cursor.Next(ctx) {
		var barber dto.Barber
		if err := cursor.Decode(&barber); err != nil {
			continue
		}
		salonName := ""
		if s, err := FetchSalonSummaryById(barber.SalonId); err == nil {
			salonName = s.Name
		}
		memberships = append(memberships, dto.SalonMembership{
			SalonId:   barber.SalonId,
			SalonName: salonName,
			Role:      dto.RoleBarber,
			BarberId:  barber.BarberId,
		})
	}
	if memberships == nil {
		memberships = []dto.SalonMembership{}
	}
	return memberships, nil
}

// FetchAllSalonsCount returns total and active salon counts.
func FetchAllSalonsCount() (total, active int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.SALONS_COLLECTION)
	n, err := col.CountDocuments(ctx, bson.M{"Deleted": false})
	if err != nil {
		return 0, 0
	}
	a, _ := col.CountDocuments(ctx, bson.M{"Deleted": false, "IsActive": true})
	return int(n), int(a)
}

// FetchAllBarbersCount returns total barber count.
func FetchAllBarbersCount() int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.BARBERS_COLLECTION)
	n, err := col.CountDocuments(ctx, bson.M{"Deleted": false})
	if err != nil {
		return 0
	}
	return int(n)
}

// FetchBarbersBySalonCount returns the number of barbers in a salon.
func FetchBarbersBySalonCount(salonId string) int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.BARBERS_COLLECTION)
	n, err := col.CountDocuments(ctx, bson.M{"SalonId": salonId, "Deleted": false})
	if err != nil {
		return 0
	}
	return int(n)
}

// FetchUserBookingsFromDB returns booking summaries for a user.
// Replaces the HTTP call to booking-service in the old user-ms.
func FetchUserBookingsFromDB(userId string) ([]dto.BookingSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.BOOKINGS_COLLECTION)
	opts := options.Find().SetSort(bson.D{{Key: "CreatedAt", Value: -1}}).SetLimit(20)
	cursor, err := col.Find(ctx, bson.M{"UserId": userId, "Deleted": false}, opts)
	if err != nil {
		return []dto.BookingSummary{}, nil
	}
	defer cursor.Close(ctx)
	var summaries []dto.BookingSummary
	for cursor.Next(ctx) {
		var b dto.Booking
		if err := cursor.Decode(&b); err != nil {
			continue
		}
		summaries = append(summaries, dto.BookingSummary{
			BookingId:    b.BookingId,
			SalonId:      b.SalonId,
			BarberId:     b.BarberId,
			Status:       b.Status,
			BookingType:  b.BookingType,
			CustomerName: b.CustomerName,
		})
	}
	if summaries == nil {
		summaries = []dto.BookingSummary{}
	}
	return summaries, nil
}

// FetchSalonBookingsFromDB returns booking summaries for a salon.
func FetchSalonBookingsFromDB(salonId string) ([]dto.BookingSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.BOOKINGS_COLLECTION)
	opts := options.Find().SetSort(bson.D{{Key: "CreatedAt", Value: -1}}).SetLimit(50)
	cursor, err := col.Find(ctx, bson.M{"SalonId": salonId, "Deleted": false}, opts)
	if err != nil {
		return []dto.BookingSummary{}, nil
	}
	defer cursor.Close(ctx)
	var summaries []dto.BookingSummary
	for cursor.Next(ctx) {
		var b dto.Booking
		if err := cursor.Decode(&b); err != nil {
			continue
		}
		summaries = append(summaries, dto.BookingSummary{
			BookingId:    b.BookingId,
			SalonId:      b.SalonId,
			BarberId:     b.BarberId,
			Status:       b.Status,
			BookingType:  b.BookingType,
			CustomerName: b.CustomerName,
		})
	}
	if summaries == nil {
		summaries = []dto.BookingSummary{}
	}
	return summaries, nil
}

// CountBookingsByStatus counts bookings with a given status from a slice.
func CountBookingsByStatus(bookings []dto.BookingSummary, status string) int {
	count := 0
	for _, b := range bookings {
		if b.Status == status {
			count++
		}
	}
	return count
}

// FetchTodayBookingsFromDB returns today's bookings for a salon.
func FetchTodayBookingsFromDB(salonId string) ([]dto.BookingSummary, error) {
	today := time.Now().Format("2006-01-02")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.TIMESLOTS_COLLECTION)
	slotCursor, err := col.Find(ctx, bson.M{"SalonId": salonId, "Date": today, "Deleted": false})
	if err != nil {
		return []dto.BookingSummary{}, nil
	}
	defer slotCursor.Close(ctx)
	var slotIds []string
	for slotCursor.Next(ctx) {
		var s dto.TimeSlot
		if err := slotCursor.Decode(&s); err == nil {
			slotIds = append(slotIds, s.SlotId)
		}
	}
	if len(slotIds) == 0 {
		return []dto.BookingSummary{}, nil
	}
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	bookingCol := dbConfig.DATABASE.Collection(dbConfig.BOOKINGS_COLLECTION)
	bCursor, err := bookingCol.Find(ctx2, bson.M{"SlotId": bson.M{"$in": slotIds}, "Deleted": false})
	if err != nil {
		log.Printf("FetchTodayBookingsFromDB: %v", err)
		return []dto.BookingSummary{}, nil
	}
	defer bCursor.Close(ctx2)
	var summaries []dto.BookingSummary
	for bCursor.Next(ctx2) {
		var b dto.Booking
		if err := bCursor.Decode(&b); err == nil {
			summaries = append(summaries, dto.BookingSummary{
				BookingId: b.BookingId, SalonId: b.SalonId, BarberId: b.BarberId,
				Status: b.Status, BookingType: b.BookingType, CustomerName: b.CustomerName,
			})
		}
	}
	if summaries == nil {
		summaries = []dto.BookingSummary{}
	}
	return summaries, nil
}
