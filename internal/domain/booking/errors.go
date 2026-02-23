package booking

import "errors"

var (
	ErrSlotAlreadyOccupied                = errors.New("slot is already occupied")
	ErrAppointmentIsFinal                 = errors.New("appointment is has final status")
	ErrAppointmentInvalidStatusTransition = errors.New("invalid status transition")
	ErrAppointmentInvalidStatus           = errors.New("invalid status")
	ErrInvalidTimeRange                   = errors.New("invalid time range")
	ErrCannotScheduleInPast               = errors.New("cannot schedule in past")
	ErrCannotBookSelf                     = errors.New("cannot book appointment to yourself")
	ErrSlotNotAvailable                   = errors.New("requested slot is not available")
)
