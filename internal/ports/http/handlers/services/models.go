package services

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/google/uuid"
)

//go:generate easyjson -all models.go

type serviceDTO struct {
	ID              uuid.UUID `json:"id"`
	PointCode       string    `json:"point_code"`
	CategoryID      int       `json:"category_id"`
	SubcategoryID   *int      `json:"subcategory_id,omitempty"`
	Name            string    `json:"name"`
	Description     *string   `json:"description,omitempty"`
	DurationMinutes int       `json:"duration_minutes"`
	Active          bool      `json:"active"`

	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
}

func toServiceDTO(s *service.Service) serviceDTO {
	var (
		minPrice float64
		maxPrice float64
	)
	if s.MinPrice != nil {
		minPrice, _ = s.MinPrice.Float64()
	}
	if s.MaxPrice != nil {
		maxPrice, _ = s.MaxPrice.Float64()
	}
	return serviceDTO{
		ID:              s.ID,
		PointCode:       s.PointCode,
		CategoryID:      s.CategoryID,
		SubcategoryID:   s.SubcategoryID,
		Name:            s.Name,
		Description:     s.Description,
		DurationMinutes: s.DurationMinutes,
		Active:          s.Active,
		MinPrice:        minPrice,
		MaxPrice:        maxPrice,
	}
}

type getServicesResponse struct {
	Services []serviceDTO `json:"services"`
}

func newGetServicesResponse(services []*service.Service) getServicesResponse {
	dtos := make([]serviceDTO, 0, len(services))
	for _, svc := range services {
		if svc.MinPrice != nil && svc.MaxPrice != nil {
			dtos = append(dtos, toServiceDTO(svc))
		}
	}
	return getServicesResponse{
		Services: dtos,
	}
}
