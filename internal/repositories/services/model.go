package services

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/google/uuid"
)

// model для SELECT с JOIN
type model struct {
	// Service fields
	ID              uuid.UUID  `db:"id"`
	PointCode       string     `db:"point_code"`
	Name            string     `db:"name"`
	Description     *string    `db:"description"`
	DurationMinutes int        `db:"duration_minutes"`
	Active          bool       `db:"active"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
	UpdatedBy       *string    `db:"updated_by"`
	Color           int        `db:"color"`

	// Category fields (from JOIN)
	CategoryID          int        `db:"category_id"`
	CategoryName        string     `db:"category_name"`
	CategoryDescription *string    `db:"category_description"`
	CategoryCreatedAt   time.Time  `db:"category_created_at"`
	CategoryUpdatedAt   *time.Time `db:"category_updated_at"`
	CategoryUpdatedBy   *string    `db:"category_updated_by"`

	// Subcategory fields (from LEFT JOIN, nullable)
	SubcategoryID          *int       `db:"subcategory_id"`
	SubcategoryName        *string    `db:"subcategory_name"`
	SubcategoryDescription *string    `db:"subcategory_description"`
	SubcategoryCreatedAt   *time.Time `db:"subcategory_created_at"`
	SubcategoryUpdatedAt   *time.Time `db:"subcategory_updated_at"`
	SubcategoryUpdatedBy   *string    `db:"subcategory_updated_by"`
}

// nolint: gosec
func (m model) convert() *service.Service {
	svc := &service.Service{
		ID:              m.ID,
		PointCode:       m.PointCode,
		Name:            m.Name,
		Description:     m.Description,
		DurationMinutes: m.DurationMinutes,
		Active:          m.Active,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		UpdatedBy:       m.UpdatedBy,
		Color:           service.Color(m.Color),
		Category: service.Category{
			ID:          m.CategoryID,
			Name:        m.CategoryName,
			Description: m.CategoryDescription,
			CreatedAt:   m.CategoryCreatedAt,
			UpdatedAt:   m.CategoryUpdatedAt,
			UpdatedBy:   m.CategoryUpdatedBy,
		},
	}

	if m.SubcategoryID != nil {
		svc.Subcategory = &service.Subcategory{
			ID:          *m.SubcategoryID,
			CategoryID:  m.CategoryID,
			Name:        *m.SubcategoryName,
			Description: m.SubcategoryDescription,
			CreatedAt:   *m.SubcategoryCreatedAt,
			UpdatedAt:   m.SubcategoryUpdatedAt,
			UpdatedBy:   m.SubcategoryUpdatedBy,
		}
	}

	return svc
}

type models []model

func (m models) convert() []*service.Service {
	list := make([]*service.Service, len(m))
	for i, item := range m {
		list[i] = item.convert()
	}
	return list
}

// writeModel для INSERT/UPDATE (только ID категорий)
type writeModel struct {
	ID              uuid.UUID
	PointCode       string
	CategoryID      int
	SubcategoryID   *int
	Name            string
	Description     *string
	DurationMinutes int
	Active          bool
	UpdatedBy       string
	Color           int
}

func convertCreateCommand(cmd *service.CreateServiceCommand) writeModel {
	return writeModel{
		PointCode:       cmd.PointCode,
		CategoryID:      cmd.CategoryID,
		SubcategoryID:   cmd.SubcategoryID,
		Name:            cmd.Name,
		Description:     cmd.Description,
		DurationMinutes: cmd.DurationMinutes,
		Active:          cmd.Active,
		UpdatedBy:       cmd.UpdatedBy,
		Color:           int(cmd.Color),
	}
}

func convertUpdateCommand(cmd *service.UpdateServiceCommand) writeModel {
	return writeModel{
		ID:              cmd.ID,
		PointCode:       cmd.PointCode,
		CategoryID:      cmd.CategoryID,
		SubcategoryID:   cmd.SubcategoryID,
		Name:            cmd.Name,
		Description:     cmd.Description,
		DurationMinutes: cmd.DurationMinutes,
		Active:          cmd.Active,
		UpdatedBy:       cmd.UpdatedBy,
		Color:           int(cmd.Color),
	}
}
