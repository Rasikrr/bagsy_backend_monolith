package booking

import "errors"

var (
	ErrAppointmentIsFinal                 = errors.New("appointment is has final status")
	ErrAppointmentInvalidStatusTransition = errors.New("invalid status transition")
	ErrAppointmentInvalidStatus           = errors.New("invalid status")
	ErrInvalidTimeRange                   = errors.New("invalid time range")
	ErrCannotScheduleInPast               = errors.New("cannot schedule in past")
)
