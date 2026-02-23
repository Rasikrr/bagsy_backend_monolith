package booking

import (
	"time"

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
	DurationMinutes int32
	MasterSlots     []MasterAvailableSlots
}

type MasterAvailableSlots struct {
	EmployeeID   uuid.UUID
	EmployeeName string
	Price        float64
	Slots        []TimeSlot
}

type TimeSlot struct {
	StartAt time.Time
	EndAt   time.Time
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
