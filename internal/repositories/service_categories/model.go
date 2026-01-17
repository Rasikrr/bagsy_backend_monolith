package servicecategories

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
)

/*
		id          SERIAL PRIMARY KEY,
	    name        TEXT NOT NULL,
	    description TEXT,
	    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
	    updated_at  TIMESTAMPTZ DEFAULT now()
	    updated_by  TEXT
*/
type model struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	UpdatedBy   *string    `db:"updated_by"`
}

func convert(e *service.Category) model {
	return model{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		UpdatedBy:   e.UpdatedBy,
	}
}

func (m model) convert() *service.Category {
	return &service.Category{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		UpdatedBy:   m.UpdatedBy,
	}
}

type models []model

func (m models) convert() []*service.Category {
	cats := make([]*service.Category, len(m))
	for i, m := range m {
		cats[i] = m.convert()
	}
	return cats
}
