package pointcategories

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
)

//go:generate easyjson -all models.go

type categoryDTO struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func toCategoryDTO(c *point.Category) categoryDTO {
	return categoryDTO{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

type getCategoriesResponse struct {
	Categories []categoryDTO `json:"categories"`
	Count      int           `json:"count"`
}

func newGetCategoriesResponse(categories []*point.Category) getCategoriesResponse {
	dtos := make([]categoryDTO, len(categories))
	for i, cat := range categories {
		dtos[i] = toCategoryDTO(cat)
	}
	return getCategoriesResponse{
		Categories: dtos,
		Count:      len(dtos),
	}
}
