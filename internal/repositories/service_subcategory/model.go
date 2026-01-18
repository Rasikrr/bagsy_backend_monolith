package servicesubcategory

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
)

/*
id          SERIAL PRIMARY KEY,
service_category_id INTEGER NOT NULL, -- FK service_categories
name        TEXT NOT NULL,
description TEXT,
created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at  TIMESTAMPTZ DEFAULT now(),
updated_by  TEXT NOT NULL DEFAULT 'system'
*/
type model struct {
	ID                int        `db:"id"`
	Name              string     `db:"name"`
	Description       *string    `db:"description"`
	ServiceCategoryID int        `db:"service_category_id"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         *time.Time `db:"updated_at"`
	UpdatedBy         *string    `db:"updated_by"`
}

func convert(e *service.Subcategory) model {
	return model{
		ID:                e.ID,
		Name:              e.Name,
		Description:       e.Description,
		ServiceCategoryID: e.CategoryID,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
		UpdatedBy:         e.UpdatedBy,
	}
}

func (m model) convert() *service.Subcategory {
	return &service.Subcategory{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CategoryID:  m.ServiceCategoryID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		UpdatedBy:   m.UpdatedBy,
	}
}

type models []model

func (m models) convert() []*service.Subcategory {
	list := make([]*service.Subcategory, len(m))
	for i, item := range m {
		list[i] = item.convert()
	}
	return list
}
