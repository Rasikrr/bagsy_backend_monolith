package masterservice

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateMasterServiceCommand struct {
	ServiceID   uuid.UUID
	Price       decimal.Decimal
	MasterPhone *string
}
