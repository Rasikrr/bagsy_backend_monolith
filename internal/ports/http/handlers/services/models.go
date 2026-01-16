package services

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
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
}

func toServiceDTO(s *entity.Service) serviceDTO {
	return serviceDTO{
		ID:              s.ID,
		PointCode:       s.PointCode,
		CategoryID:      s.CategoryID,
		SubcategoryID:   s.SubcategoryID,
		Name:            s.Name,
		Description:     s.Description,
		DurationMinutes: s.DurationMinutes,
		Active:          s.Active,
	}
}

type getServicesResponse struct {
	Services []serviceDTO `json:"services"`
}

func newGetServicesResponse(services []*entity.Service) getServicesResponse {
	dtos := make([]serviceDTO, 0, len(services))
	for _, svc := range services {
		dtos = append(dtos, toServiceDTO(svc))
	}
	return getServicesResponse{
		Services: dtos,
	}
}
