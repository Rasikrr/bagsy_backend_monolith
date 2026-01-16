package pointcategoryservices

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
)

/*
id                   SERIAL PRIMARY KEY,
point_category_id    INTEGER NOT NULL,
service_category_id  INTEGER NOT NULL,
created_at           TIMESTAMPTZ NOT NULL DEFAULT now()
*/
type model struct {
	ID                int       `db:"id"`
	PointCategoryID   int       `db:"point_category_id"`
	ServiceCategoryID int       `db:"service_category_id"`
	CreatedAt         time.Time `db:"created_at"`
}

func convert(e *point.CategoryService) model {
	return model{
		ID:                e.ID,
		PointCategoryID:   e.PointCategoryID,
		ServiceCategoryID: e.ServiceCategoryID,
		CreatedAt:         e.CreatedAt,
	}
}

func (m model) convert() *point.CategoryService {
	return &point.CategoryService{
		ID:                m.ID,
		PointCategoryID:   m.PointCategoryID,
		ServiceCategoryID: m.ServiceCategoryID,
		CreatedAt:         m.CreatedAt,
	}
}

type models []model

func (m models) convert() []*point.CategoryService {
	list := make([]*point.CategoryService, len(m))
	for i, item := range m {
		list[i] = item.convert()
	}
	return list
}
