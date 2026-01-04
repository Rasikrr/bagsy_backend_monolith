package bagsies

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type model struct {
	ID           uuid.UUID       `db:"id"`
	PointCode    string          `db:"point_code"`
	ClientPhone  string          `db:"client_phone"`
	Status       string          `db:"status"`
	MasterPhone  string          `db:"master_phone"`
	Price        decimal.Decimal `db:"price"`
	ServiceID    uuid.UUID       `db:"service_id"`
	StartAt      time.Time       `db:"start_at"`
	EndAt        time.Time       `db:"end_at"`
	Comment      *string         `db:"comment"`
	RejectReason *string         `db:"reject_reason"`
	CreatedAt    time.Time       `db:"created_at"`
	UpdatedAt    *time.Time      `db:"updated_at"`
	UpdatedBy    string          `db:"updated_by"`
}

type models []model

func (m model) convert() (*entity.Bagsy, error) {
	status, err := enum.BagsyStatusString(m.Status)
	if err != nil {
		return nil, err
	}
	bagsy := &entity.Bagsy{
		ID:           m.ID,
		PointCode:    m.PointCode,
		MasterPhone:  m.MasterPhone,
		Status:       status,
		ServiceID:    m.ServiceID,
		ClientPhone:  m.ClientPhone,
		Price:        m.Price,
		StartAt:      m.StartAt,
		EndAt:        m.EndAt,
		Comment:      m.Comment,
		RejectReason: m.RejectReason,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		UpdatedBy:    m.UpdatedBy,
	}
	return bagsy, nil
}

func (ms models) convert() ([]*entity.Bagsy, error) {
	res := make([]*entity.Bagsy, len(ms))
	for i, m := range ms {
		out, err := m.convert()
		if err != nil {
			return nil, err
		}
		res[i] = out
	}
	return res, nil
}

func convertToModel(b *entity.Bagsy) *model {
	return &model{
		ID:           b.ID,
		Status:       b.Status.String(),
		PointCode:    b.PointCode,
		ClientPhone:  b.ClientPhone,
		Price:        b.Price,
		ServiceID:    b.ServiceID,
		MasterPhone:  b.MasterPhone,
		StartAt:      b.StartAt,
		EndAt:        b.EndAt,
		Comment:      b.Comment,
		RejectReason: b.RejectReason,
		CreatedAt:    b.CreatedAt,
		UpdatedAt:    b.UpdatedAt,
		UpdatedBy:    b.UpdatedBy,
	}
}
