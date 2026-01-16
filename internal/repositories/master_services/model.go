package masterservices

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/master_service"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

/*
id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
master_phone TEXT NOT NULL,
service_id   UUID NOT NULL,
price        DECIMAL(10,2) NOT NULL,
active       BOOLEAN DEFAULT false,
created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at   TIMESTAMPTZ DEFAULT now()
*/
type model struct {
	ID          uuid.UUID       `db:"id"`
	MasterPhone string          `db:"master_phone"`
	ServiceID   uuid.UUID       `db:"service_id"`
	Price       decimal.Decimal `db:"price"`
	Active      bool            `db:"active"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   *time.Time      `db:"updated_at"`
	UpdatedBy   *string         `db:"updated_by"`
}

func convert(e *masterservice.MasterService) model {
	return model{
		ID:          e.ID,
		MasterPhone: e.MasterPhone,
		ServiceID:   e.ServiceID,
		Price:       e.Price,
		Active:      e.Active,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		UpdatedBy:   e.UpdatedBy,
	}
}

func (m model) convert() *masterservice.MasterService {
	return &masterservice.MasterService{
		ID:          m.ID,
		MasterPhone: m.MasterPhone,
		ServiceID:   m.ServiceID,
		Price:       m.Price,
		Active:      m.Active,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		UpdatedBy:   m.UpdatedBy,
	}
}

type models []model

func (m models) convert() []*masterservice.MasterService {
	list := make([]*masterservice.MasterService, len(m))
	for i, item := range m {
		list[i] = item.convert()
	}
	return list
}
