package booking

import (
	"time"

	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────────
// Get Available Slots
// ─────────────────────────────────────────────────────────────────

type getSlotsRequest struct {
	LocationID uuid.UUID  `json:"location_id"`
	ServiceID  uuid.UUID  `json:"service_id"`
	EmployeeID *uuid.UUID `json:"employee_id"` // optional
	StartDate  time.Time  `json:"start_date"`
	EndDate    time.Time  `json:"end_date"`
}

type getSlotsResponse struct {
	ServiceID       uuid.UUID        `json:"service_id"`
	LocationID      uuid.UUID        `json:"location_id"`
	DurationMinutes int              `json:"duration_minutes"`
	MasterSlots     []masterTimeSlot `json:"master_slots"`
}

type masterTimeSlot struct {
	EmployeeID   uuid.UUID  `json:"employee_id"`
	EmployeeName string     `json:"employee_name"`
	Price        float64    `json:"price"`
	Slots        []timeSlot `json:"slots"`
}

type timeSlot struct {
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}

// ─────────────────────────────────────────────────────────────────
// Create Booking
// ─────────────────────────────────────────────────────────────────

type createRequest struct {
	LocationID uuid.UUID `json:"location_id"`
	ServiceID  uuid.UUID `json:"service_id"`
	EmployeeID uuid.UUID `json:"employee_id"`
	StartAt    time.Time `json:"start_at"`

	// Customer Info
	Phone     string  `json:"phone"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
	Comment   *string `json:"comment"`
}

type createResponse struct {
	ID uuid.UUID `json:"id"`
}

// ─────────────────────────────────────────────────────────────────
// Confirm Booking
// ─────────────────────────────────────────────────────────────────

type confirmRequest struct {
	Code string `json:"code"`
}

// ─────────────────────────────────────────────────────────────────
// Cancel Booking
// ─────────────────────────────────────────────────────────────────

type cancelRequest struct {
	Reason string `json:"reason"`
}
