package booking

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

// CalendarEntry — read-only projection for calendar listing.
// Assembled via JOIN query, not a domain aggregate.
type CalendarEntry struct {
	AppointmentID   uuid.UUID
	Status          Status
	StartAt         time.Time
	EndAt           time.Time
	Price           shared.Money
	DurationMinutes shared.Duration
	CustomerComment *string

	EmployeeID   uuid.UUID
	EmployeeName string

	CustomerID    uuid.UUID
	CustomerName  string
	CustomerPhone shared.Phone

	ServiceID    uuid.UUID
	ServiceName  string
	ServiceColor catalog.Color

	LocationID   uuid.UUID
	LocationName string
}
