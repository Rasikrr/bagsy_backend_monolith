package entity

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Bagsy struct {
	ID           uuid.UUID
	ServiceID    uuid.UUID
	PointCode    string
	ClientPhone  string
	MasterPhone  string
	Status       enum.BagsyStatus
	Price        decimal.Decimal
	StartAt      time.Time
	EndAt        time.Time
	Comment      *string
	RejectReason *string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	UpdatedBy    string
}
