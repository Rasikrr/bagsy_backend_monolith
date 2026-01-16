package services

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/go-chi/chi/v5"
)

type servicesService interface {
	GetByPointCode(ctx context.Context, pointCode string) ([]*service.Service, error)
}

type Controller struct {
	servicesService servicesService
}

func New(servicesService servicesService) *Controller {
	return &Controller{
		servicesService: servicesService,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	router.Route("/api/v1/services", func(r chi.Router) {
		r.Get("/{point_code}", c.getServicesByPointCode)
	})
}
