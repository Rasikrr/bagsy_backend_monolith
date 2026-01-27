package services

import (
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

//go:generate easyjson -all models.go

type createServiceRequest struct {
	PointCode       string  `json:"point_code" validate:"required"`
	CategoryID      int     `json:"category_id" validate:"required,gt=0"`
	SubcategoryID   *int    `json:"subcategory_id,omitempty"`
	Name            string  `json:"name" validate:"required,min=1,max=255"`
	Description     *string `json:"description,omitempty"`
	DurationMinutes int     `json:"duration_minutes" validate:"required,gt=0"`
	Color           string  `json:"color" validate:"required" enums:"blue,green,red,yellow,purple,orange,gray"`
}

type createServiceResponse struct {
	ServiceID uuid.UUID `json:"service_id"`
}

func newCreateServiceResponse(serviceID uuid.UUID) *createServiceResponse {
	return &createServiceResponse{
		ServiceID: serviceID,
	}
}

func (r *createServiceRequest) Validate() error {
	err := request.GetValidator().Struct(r)
	if err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

func (r *createServiceRequest) toCommand(updatedBy string) (*service.CreateServiceCommand, error) {
	color, err := service.ColorString(r.Color)
	if err != nil {
		return nil, domainErr.NewInvalidInputError("invalid color value", err)
	}
	return &service.CreateServiceCommand{
		PointCode:       r.PointCode,
		CategoryID:      r.CategoryID,
		SubcategoryID:   r.SubcategoryID,
		Name:            r.Name,
		Description:     r.Description,
		DurationMinutes: r.DurationMinutes,
		UpdatedBy:       updatedBy,
		Color:           color,
	}, nil
}

type serviceDTO struct {
	ID              uuid.UUID `json:"id"`
	PointCode       string    `json:"point_code"`
	CategoryID      int       `json:"category_id"`
	SubcategoryID   *int      `json:"subcategory_id,omitempty"`
	Name            string    `json:"name"`
	Description     *string   `json:"description,omitempty"`
	DurationMinutes int       `json:"duration_minutes"`
	Active          bool      `json:"active"`
	Color           string    `json:"color"`

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
		Color:           s.Color.String(),
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
		dtos = append(dtos, toServiceDTO(svc))
	}
	return getServicesResponse{
		Services: dtos,
	}
}

type getServicesRequest struct {
	Active *bool `query:"is_active"`
}

func (g *getServicesRequest) GetQueryParameters(r *http.Request) error {
	active := r.URL.Query().Get("is_active")

	if active != "" {
		activeBool, err := strconv.ParseBool(active)
		if err != nil {
			return request.HandleValidationError(err)
		}
		g.Active = &activeBool
	}

	return nil
}
