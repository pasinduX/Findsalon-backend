package service

import (
	"context"
	"time"

	"findsalon-backend/dto"
	salonerr "findsalon-backend/errors"
	"findsalon-backend/repository"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingService struct {
	bookingRepo  *repository.BookingRepository
	scheduleRepo *repository.ScheduleRepository
	blockRepo    *repository.BlockRepository
	serviceRepo  *repository.ServiceRepository
}

func NewBookingService(
	bookingRepo *repository.BookingRepository,
	scheduleRepo *repository.ScheduleRepository,
	blockRepo *repository.BlockRepository,
	serviceRepo *repository.ServiceRepository,
) *BookingService {
	return &BookingService{bookingRepo, scheduleRepo, blockRepo, serviceRepo}
}

func (s *BookingService) GetAvailability(ctx context.Context, req dto.AvailabilityRequest) ([]dto.AvailableSlot, error) {
	svc, err := s.serviceRepo.FindById(ctx, req.ServiceId)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		return nil, salonerr.ErrNotFound
	}
	schedule, err := s.scheduleRepo.FindByBarberId(ctx, req.BarberId)
	if err != nil {
		return nil, err
	}
	if schedule == nil {
		return nil, salonerr.ErrNotFound
	}
	window, err := GetWorkWindow(*schedule, req.Date)
	if err != nil {
		return nil, err
	}
	bookings, err := s.bookingRepo.FindActiveByBarberAndWindow(ctx, req.BarberId, window.Start, window.End)
	if err != nil {
		return nil, err
	}
	blocks, err := s.blockRepo.FindByBarberAndWindow(ctx, req.BarberId, window.Start, window.End)
	if err != nil {
		return nil, err
	}
	return GetAvailableSlots(AvailabilityInput{
		Window:   window,
		Service:  *svc,
		StepMin:  schedule.EffectiveSlotStep(),
		Bookings: bookings,
		Blocks:   blocks,
		Timezone: schedule.EffectiveTimezone(),
	})
}

func (s *BookingService) CreateDirectBooking(ctx context.Context, req dto.DirectBookingRequest) (*dto.Booking, error) {
	svc, err := s.serviceRepo.FindById(ctx, req.ServiceId)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		return nil, salonerr.ErrNotFound
	}
	endTime := req.StartTime.UTC().Add(svc.TotalDuration())
	if req.StartTime.UTC().Before(time.Now().UTC()) {
		return nil, salonerr.ErrPastTime
	}
	bookings, err := s.bookingRepo.FindActiveByBarberAndWindow(ctx, req.BarberId, req.StartTime, endTime)
	if err != nil {
		return nil, err
	}
	blocks, err := s.blockRepo.FindByBarberAndWindow(ctx, req.BarberId, req.StartTime, endTime)
	if err != nil {
		return nil, err
	}
	if !isAvailable(req.StartTime.UTC(), endTime, bookings, blocks) {
		return nil, salonerr.ErrSlotTaken
	}
	now := time.Now().UTC()
	booking := dto.Booking{
		BookingId:    uuid.New().String(),
		BarberId:     req.BarberId,
		SalonId:      req.SalonId,
		ServiceId:    req.ServiceId,
		UserId:       req.CustomerId,
		CustomerName: req.CustomerName,
		Notes:        req.Notes,
		StartTime:    req.StartTime.UTC(),
		EndTime:      endTime,
		Status:       dto.BookingStatusBooked,
		BookingType:  dto.BookingTypeDirect,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.bookingRepo.Insert(ctx, booking); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, salonerr.ErrSlotTaken
		}
		return nil, err
	}
	return &booking, nil
}
