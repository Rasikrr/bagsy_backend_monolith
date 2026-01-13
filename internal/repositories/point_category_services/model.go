package pointcategoryservices

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
)

type model struct {
	ID                int       `db:"id"`
	PointCategoryID   int       `db:"point_category_id"`
	ServiceCategoryID int       `db:"service_category_id"`
	CreatedAt         time.Time `db:"created_at"`
}

func (m model) convert() *entity.PointCategoryService {
	return &entity.PointCategoryService{
		ID:                m.ID,
		PointCategoryID:   m.PointCategoryID,
		ServiceCategoryID: m.ServiceCategoryID,
		CreatedAt:         m.CreatedAt,
	}
}

func convert(e *entity.PointCategoryService) model {
	return model{
		ID:                e.ID,
		PointCategoryID:   e.PointCategoryID,
		ServiceCategoryID: e.ServiceCategoryID,
		CreatedAt:         e.CreatedAt,
	}
}

type models []model

func (m models) convert() []*entity.PointCategoryService {
	out := make([]*entity.PointCategoryService, len(m))
	for i := range m {
		out[i] = m[i].convert()
	}
	return out
}

// serviceCategoryModel для JOIN запросов
type serviceCategoryModel struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	UpdatedBy   *string    `db:"updated_by"`
}

func (m serviceCategoryModel) convert() *entity.ServiceCategory {
	return &entity.ServiceCategory{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		UpdatedBy:   m.UpdatedBy,
	}
}

type serviceCategoryModels []serviceCategoryModel

func (m serviceCategoryModels) convert() []*entity.ServiceCategory {
	out := make([]*entity.ServiceCategory, len(m))
	for i := range m {
		out[i] = m[i].convert()
	}
	return out
}

// pointCategoryModel для JOIN запросов
type pointCategoryModel struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	UpdatedBy   *string    `db:"updated_by"`
}

func (m pointCategoryModel) convert() *entity.PointCategory {
	return &entity.PointCategory{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		UpdatedBy:   m.UpdatedBy,
	}
}

type pointCategoryModels []pointCategoryModel

func (m pointCategoryModels) convert() []*entity.PointCategory {
	out := make([]*entity.PointCategory, len(m))
	for i := range m {
		out[i] = m[i].convert()
	}
	return out
}
