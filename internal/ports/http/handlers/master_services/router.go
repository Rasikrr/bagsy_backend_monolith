package masterservices

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/master_service"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/go-chi/chi/v5"
)

type masterServicesService interface {
	Create(ctx context.Context, cmd *masterservice.CreateMasterServiceCommand) (*masterservice.MasterService, error)
}

type Controller struct {
	masterServicesService masterServicesService
	authMiddleware        *middlewares.Auth
}

func New(
	masterServicesService masterServicesService,
	authMiddleware *middlewares.Auth,
) *Controller {
	return &Controller{
		masterServicesService: masterServicesService,
		authMiddleware:        authMiddleware,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	workers := c.authMiddleware.AuthorizeWorkers()
	router.Route("/api/v1/master-services", func(r chi.Router) {
		r.With(workers).Post("/", c.createMasterService)
	})
}
