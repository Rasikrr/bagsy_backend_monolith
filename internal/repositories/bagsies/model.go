package bagsies

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
)

type model struct {
	ID        string    `db:"id"`
	PointCode string    `db:"point_code"`
	Phone     string    `db:"phone"`
	StartAt   time.Time `db:"start_at"`
	EndAt     time.Time `db:"end_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	UpdatedBy *string   `db:"updated_by"`
}

type models []model

func (m model) convert() *entity.Bagsy {
	return &entity.Bagsy{
		ID:        m.ID,
		PointCode: m.PointCode,
		UserPhone: m.Phone,
		StartAt:   m.StartAt,
		EndAt:     m.EndAt,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		UpdatedBy: m.UpdatedBy,
	}
}

func (ms models) convert() []*entity.Bagsy {
	res := make([]*entity.Bagsy, len(ms))
	for i, m := range ms {
		res[i] = m.convert()
	}
	return res
}

func convertToModel(b *entity.Bagsy) *model {
	return &model{
		ID:        b.ID,
		PointCode: b.PointCode,
		Phone:     b.UserPhone,
		StartAt:   b.StartAt,
		EndAt:     b.EndAt,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
		UpdatedBy: b.UpdatedBy,
	}
}
