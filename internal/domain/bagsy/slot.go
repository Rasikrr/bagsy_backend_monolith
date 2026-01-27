package bagsy

import (
	"github.com/shopspring/decimal"
	"time"

	"github.com/google/uuid"
)

// TimeSlot represents an available booking slot
type TimeSlot struct {
	StartAt time.Time
	EndAt   time.Time
}

// MasterSlot represents available slots for a specific master
type MasterSlot struct {
	MasterPhone        string
	MasterName         string
	MasterServicePrice decimal.Decimal
	Slots              []TimeSlot
}

// AvailableSlots represents the response for available slots query
type AvailableSlots struct {
	ServiceID       uuid.UUID
	PointCode       string
	DurationMinutes int
	MasterSlots     []MasterSlot
}
