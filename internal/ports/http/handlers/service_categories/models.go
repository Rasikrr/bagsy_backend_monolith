package servicecategories

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
)

//go:generate easyjson -all models.go

type categoryDTO struct {
	ID            int              `json:"id"`
	Name          string           `json:"name"`
	Description   *string          `json:"description,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     *time.Time       `json:"updated_at,omitempty"`
	Subcategories []subcategoryDTO `json:"subcategories"`
}

type subcategoryDTO struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func toCategoryDTO(cat *service.CategoryWithSubcategories) categoryDTO {
	subDTOs := make([]subcategoryDTO, len(cat.Subcategories))
	for i, sub := range cat.Subcategories {
		subDTOs[i] = subcategoryDTO{
			ID:          sub.ID,
			Name:        sub.Name,
			Description: sub.Description,
			CreatedAt:   sub.CreatedAt,
			UpdatedAt:   sub.UpdatedAt,
		}
	}
	return categoryDTO{
		ID:            cat.Category.ID,
		Name:          cat.Category.Name,
		Description:   cat.Category.Description,
		CreatedAt:     cat.Category.CreatedAt,
		UpdatedAt:     cat.Category.UpdatedAt,
		Subcategories: subDTOs,
	}
}

type getByPointCodeResponse struct {
	Categories []categoryDTO `json:"categories"`
	Count      int           `json:"count"`
}

func newGetByPointCodeResponse(categories []*service.CategoryWithSubcategories) getByPointCodeResponse {
	dtos := make([]categoryDTO, len(categories))
	for i, cat := range categories {
		dtos[i] = toCategoryDTO(cat)
	}
	return getByPointCodeResponse{
		Categories: dtos,
		Count:      len(dtos),
	}
}
