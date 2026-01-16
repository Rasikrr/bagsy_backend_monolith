package pointcategories

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
)

type model struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	UpdatedBy   *string    `db:"updated_by"`
}

func convert(e *point.Category) model {
	return model{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		UpdatedBy:   e.UpdatedBy,
	}
}

func (m model) convert() *point.Category {
	return &point.Category{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		UpdatedBy:   m.UpdatedBy,
	}
}

type models []model

func (m models) convert() []*point.Category {
	out := make([]*point.Category, len(m))
	for i := range m {
		out[i] = m[i].convert()
	}
	return out
}
