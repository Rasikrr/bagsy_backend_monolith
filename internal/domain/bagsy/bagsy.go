package bagsy

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Bagsy struct {
	ID           uuid.UUID
	ServiceID    uuid.UUID
	PointCode    string
	ClientPhone  string
	MasterPhone  string
	Status       Status
	Price        decimal.Decimal
	StartAt      time.Time
	EndAt        time.Time
	Comment      *string
	RejectReason *string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	UpdatedBy    string
}
