package errors

import "errors"

var (
	ErrSlotTaken        = errors.New("time slot already booked")
	ErrNotWorkingDay    = errors.New("barber is not working on this day")
	ErrOutsideWorkHours = errors.New("requested time is outside working hours")
	ErrPastTime         = errors.New("cannot book a slot in the past")
	ErrNotFound         = errors.New("resource not found")
	ErrInvalidInput     = errors.New("invalid input")
	ErrInvalidTimezone  = errors.New("invalid or unknown timezone")
)
