package repository

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ServiceRepository reads SalonService records for the availability engine.
type ServiceRepository struct{}

func (r *ServiceRepository) FindById(ctx context.Context, serviceId string) (*dto.SalonService, error) {
	col := dbConfig.DATABASE.Collection(dbConfig.SALONSERVICES_COLLECTION)
	var s dto.SalonService
	err := col.FindOne(ctx, bson.M{"ServiceId": serviceId}).Decode(&s)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &s, err
}

// ScheduleRepository reads WeeklySchedule records.
type ScheduleRepository struct{}

func (r *ScheduleRepository) FindByBarberId(ctx context.Context, barberId string) (*dto.WeeklySchedule, error) {
	col := dbConfig.DATABASE.Collection(dbConfig.WEEKLY_SCHEDULES_COLLECTION)
	var s dto.WeeklySchedule
	err := col.FindOne(ctx, bson.M{"BarberId": barberId}).Decode(&s)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &s, err
}

// BlockRepository handles ScheduleBlock persistence.
type BlockRepository struct{}

func (r *BlockRepository) FindByBarberAndWindow(ctx context.Context, barberId string, windowStart, windowEnd time.Time) ([]dto.ScheduleBlock, error) {
	col := dbConfig.DATABASE.Collection(dbConfig.SCHEDULE_BLOCKS_COLLECTION)
	filter := bson.M{
		"BarberId":  barberId,
		"StartTime": bson.M{"$lt": windowEnd},
		"EndTime":   bson.M{"$gt": windowStart},
	}
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var blocks []dto.ScheduleBlock
	return blocks, cursor.All(ctx, &blocks)
}

func (r *BlockRepository) FindByBarberDate(ctx context.Context, barberId, date string) ([]dto.ScheduleBlock, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	dayStart := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return r.FindByBarberAndWindow(ctx, barberId, dayStart, dayStart.Add(24*time.Hour))
}

func (r *BlockRepository) Insert(ctx context.Context, block dto.ScheduleBlock) error {
	col := dbConfig.DATABASE.Collection(dbConfig.SCHEDULE_BLOCKS_COLLECTION)
	_, err := col.InsertOne(ctx, block)
	return err
}

func (r *BlockRepository) FindById(ctx context.Context, blockId string) (*dto.ScheduleBlock, error) {
	col := dbConfig.DATABASE.Collection(dbConfig.SCHEDULE_BLOCKS_COLLECTION)
	var b dto.ScheduleBlock
	err := col.FindOne(ctx, bson.M{"BlockId": blockId}).Decode(&b)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &b, err
}

func (r *BlockRepository) Update(ctx context.Context, block dto.ScheduleBlock) error {
	col := dbConfig.DATABASE.Collection(dbConfig.SCHEDULE_BLOCKS_COLLECTION)
	_, err := col.UpdateOne(ctx, bson.M{"BlockId": block.BlockId}, bson.M{"$set": block})
	return err
}

func (r *BlockRepository) Delete(ctx context.Context, blockId string) error {
	col := dbConfig.DATABASE.Collection(dbConfig.SCHEDULE_BLOCKS_COLLECTION)
	_, err := col.DeleteOne(ctx, bson.M{"BlockId": blockId})
	return err
}

// BookingRepository handles Booking persistence for the availability engine.
type BookingRepository struct{}

func (r *BookingRepository) FindActiveByBarberAndWindow(ctx context.Context, barberId string, windowStart, windowEnd time.Time) ([]dto.Booking, error) {
	col := dbConfig.DATABASE.Collection(dbConfig.BOOKINGS_COLLECTION)
	filter := bson.M{
		"BarberId":  barberId,
		"Deleted":   false,
		"Status":    bson.M{"$nin": []string{dto.BookingStatusCancelled}},
		"StartTime": bson.M{"$lt": windowEnd},
		"EndTime":   bson.M{"$gt": windowStart},
	}
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var bookings []dto.Booking
	return bookings, cursor.All(ctx, &bookings)
}

func (r *BookingRepository) Insert(ctx context.Context, booking dto.Booking) error {
	col := dbConfig.DATABASE.Collection(dbConfig.BOOKINGS_COLLECTION)
	_, err := col.InsertOne(ctx, booking)
	return err
}

func (r *BookingRepository) CancelById(ctx context.Context, bookingId string) error {
	col := dbConfig.DATABASE.Collection(dbConfig.BOOKINGS_COLLECTION)
	_, err := col.UpdateOne(ctx,
		bson.M{"BookingId": bookingId},
		bson.M{"$set": bson.M{"Status": dto.BookingStatusCancelled, "UpdatedAt": time.Now().UTC()}},
	)
	return err
}
