package services

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/google/uuid"
)

//go:generate easyjson -all models.go

type categoryDTO struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type subcategoryDTO struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type serviceDTO struct {
	ID              uuid.UUID       `json:"id"`
	PointCode       string          `json:"point_code"`
	Category        categoryDTO     `json:"category"`
	Subcategory     *subcategoryDTO `json:"subcategory,omitempty"`
	Name            string          `json:"name"`
	Description     *string         `json:"description,omitempty"`
	DurationMinutes int             `json:"duration_minutes"`
	Active          bool            `json:"active"`
	Color           string          `json:"color"`
}

func toServiceDTO(s *service.Service) serviceDTO {
	dto := serviceDTO{
		ID:        s.ID,
		PointCode: s.PointCode,
		Category: categoryDTO{
			ID:          s.Category.ID,
			Name:        s.Category.Name,
			Description: s.Category.Description,
		},
		Name:            s.Name,
		Description:     s.Description,
		DurationMinutes: s.DurationMinutes,
		Active:          s.Active,
		Color:           s.Color.String(),
	}

	if s.Subcategory != nil {
		dto.Subcategory = &subcategoryDTO{
			ID:          s.Subcategory.ID,
			Name:        s.Subcategory.Name,
			Description: s.Subcategory.Description,
		}
	}

	return dto
}

type getServicesResponse struct {
	Services []serviceDTO `json:"services"`
}

func newGetServicesResponse(services []*service.Service) getServicesResponse {
	dtos := make([]serviceDTO, 0, len(services))
	for _, svc := range services {
		dtos = append(dtos, toServiceDTO(svc))
	}
	return getServicesResponse{
		Services: dtos,
	}
}
