package bagsies

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/google/uuid"
)

type model struct {
	ID          uuid.UUID  `db:"id"`
	PointCode   string     `db:"point_code"`
	UserPhone   string     `db:"user_phone"`
	Status      string     `db:"status"`
	MasterPhone string     `db:"master_phone"`
	ServiceID   uuid.UUID  `db:"service_id"`
	StartAt     time.Time  `db:"start_at"`
	EndAt       time.Time  `db:"end_at"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	UpdatedBy   *string    `db:"updated_by"`
}

type models []model

func (m model) convert() (*entity.Bagsy, error) {
	status, err := enum.BagsyStatusString(m.Status)
	if err != nil {
		return nil, err
	}
	bagsy := &entity.Bagsy{
		ID:          m.ID,
		PointCode:   m.PointCode,
		MasterPhone: m.MasterPhone,
		Status:      status,
		ServiceID:   m.ServiceID,
		UserPhone:   m.UserPhone,
		StartAt:     m.StartAt,
		EndAt:       m.EndAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		UpdatedBy:   m.UpdatedBy,
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
		ID:          b.ID,
		Status:      b.Status.String(),
		PointCode:   b.PointCode,
		UserPhone:   b.UserPhone,
		ServiceID:   b.ServiceID,
		MasterPhone: b.MasterPhone,
		StartAt:     b.StartAt,
		EndAt:       b.EndAt,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
		UpdatedBy:   b.UpdatedBy,
	}
}
