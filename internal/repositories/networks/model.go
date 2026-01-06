package networks

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
)

type model struct {
	Code        string     `db:"code"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	CreatedBy   string     `db:"created_by"`
	UpdatedBy   string     `db:"updated_by"`
}

func convert(e *entity.Network) model {
	return model{
		Code:        e.Code,
		Name:        e.Name,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		DeletedAt:   e.DeletedAt,
		CreatedBy:   e.CreatedBy,
		UpdatedBy:   e.UpdatedBy,
	}
}

func (m model) convert() *entity.Network {
	return &entity.Network{
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
	}
}

type models []model

func (m models) convert() []*entity.Network {
	out := make([]*entity.Network, len(m))
	for i, m := range m {
		out[i] = m.convert()
	}
	return out
}
