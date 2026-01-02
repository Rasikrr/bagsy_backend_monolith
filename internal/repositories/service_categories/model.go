package servicecategories

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
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

func convert(e *entity.ServiceCategory) model {
	return model{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		UpdatedBy:   e.UpdatedBy,
	}
}

func (m model) convert() *entity.ServiceCategory {
	return &entity.ServiceCategory{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		UpdatedBy:   m.UpdatedBy,
	}
}

type models []model

func (m models) convert() []*entity.ServiceCategory {
	cats := make([]*entity.ServiceCategory, len(m))
	for i, m := range m {
		cats[i] = m.convert()
	}
	return cats
}
