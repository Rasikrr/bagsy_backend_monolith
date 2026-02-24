package booking

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/booking"
	"github.com/google/uuid"
)

type GetAvailableSlotsInput struct {
	LocationID uuid.UUID
	ServiceID  uuid.UUID
	EmployeeID *uuid.UUID
	StartDate  time.Time
	EndDate    time.Time
}

type GetAvailableSlotsOutput struct {
	ServiceID       uuid.UUID
	LocationID      uuid.UUID
	DurationMinutes int
	MasterSlots     []MasterAvailableSlots
}

type MasterAvailableSlots struct {
	EmployeeID   uuid.UUID
	EmployeeName string
	Price        float64
	Slots        []booking.TimeSlot
}

type CreateBookingInput struct {
	LocationID uuid.UUID
	ServiceID  uuid.UUID
	EmployeeID uuid.UUID
	StartAt    time.Time

	Phone     string
	FirstName string
	LastName  *string
	Comment   *string
}

type CreateBookingOutput struct {
	ID uuid.UUID
}

// Calendar

type GetCalendarInput struct {
	LocationID       *uuid.UUID
	EmployeeID       *uuid.UUID
	StartDate        time.Time
	EndDate          time.Time
	IncludeCancelled bool
}
