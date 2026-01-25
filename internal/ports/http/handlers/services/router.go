package services

import (
	"context"
	"github.com/google/uuid"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/go-chi/chi/v5"
)

type servicesService interface {
	GetByPointCode(ctx context.Context, pointCode string, isActive *bool) ([]*service.Service, error)
	Create(ctx context.Context, cmd *service.CreateServiceCommand) (uuid.UUID, error)
}

type Controller struct {
	servicesService servicesService
	authMiddleware  *middlewares.Auth
}

func New(servicesService servicesService, authMiddleware *middlewares.Auth) *Controller {
	return &Controller{
		servicesService: servicesService,
		authMiddleware:  authMiddleware,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	management := c.authMiddleware.AuthorizeManagement()
	router.Route("/api/v1/services", func(r chi.Router) {
		r.Get("/{point_code}", c.getServicesByPointCode)

		r.With(management).Post("/", c.createService)
	})
}
