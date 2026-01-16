package masterservice

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type MasterService struct {
	ID          uuid.UUID
	MasterPhone string
	ServiceID   uuid.UUID
	Price       decimal.Decimal
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	UpdatedBy   *string
}
