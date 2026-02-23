package booking

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CalendarEntry — read-only projection for calendar listing.
// Assembled via JOIN query, not a domain aggregate.
type CalendarEntry struct {
	AppointmentID   uuid.UUID
	Status          Status
	StartAt         time.Time
	EndAt           time.Time
	Price           decimal.Decimal
	DurationMinutes int
	CustomerComment *string

	EmployeeID   uuid.UUID
	EmployeeName string

	CustomerID    uuid.UUID
	CustomerName  string
	CustomerPhone string

	ServiceID    uuid.UUID
	ServiceName  string
	ServiceColor string

	LocationID   uuid.UUID
	LocationName string
}
