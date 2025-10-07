package bagsies

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
)

type model struct {
	ID            string    `db:"id"`
	PointCode     string    `db:"point_code"`
	ProviderPhone string    `db:"provider_phone"`
	UserPhone     string    `db:"user_phone"`
	FirstName     string    `db:"first_name"`
	LastName      string    `db:"last_name"`
	Description   string    `db:"description"`
	Service       string    `db:"service"`
	StartAt       time.Time `db:"start_at"`
	EndAt         time.Time `db:"end_at"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	UpdatedBy     *string   `db:"updated_by"`
}

type models []model

func (m model) convert() *entity.Bagsy {
	bagsy := &entity.Bagsy{
		ID:            m.ID,
		PointCode:     m.PointCode,
		ProviderPhone: m.ProviderPhone,
		UserPhone:     m.UserPhone,
		FirstName:     m.FirstName,
		LastName:      m.LastName,
		Description:   m.Description,
		Service:       m.Service,
		StartAt:       m.StartAt,
		EndAt:         m.EndAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}

	if m.UpdatedBy != nil {
		bagsy.UpdatedBy = *m.UpdatedBy
	}

	return bagsy
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
		ID:            b.ID,
		PointCode:     b.PointCode,
		ProviderPhone: b.ProviderPhone,
		UserPhone:     b.UserPhone,
		FirstName:     b.FirstName,
		LastName:      b.LastName,
		Description:   b.Description,
		Service:       b.Service,
		StartAt:       b.StartAt,
		EndAt:         b.EndAt,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
		UpdatedBy:     &b.UpdatedBy,
	}
}
