package bagsy

import (
	"time"

	"github.com/google/uuid"
)

// CreateBagsyCommand - команда для создания брони
type CreateBagsyCommand struct {
	ServiceID   uuid.UUID
	MasterPhone string

	StartAt time.Time

	ClientPhone string
	Name        string
	Surname     string
	Comment     *string
}

// GetAvailableSlotsCommand - команда для получения свободных слотов
type GetAvailableSlotsCommand struct {
	PointCode   string
	ServiceID   uuid.UUID
	MasterPhone *string   // optional filter by specific master
	StartDate   time.Time // start of period (default: now)
	EndDate     time.Time // end of period (default: now + 2 weeks)
}
