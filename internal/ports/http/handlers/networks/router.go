package networks

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/command"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/go-chi/chi/v5"
)

type networksService interface {
	Create(ctx context.Context, createReq *command.CreateNetworkCommand) error
}
type Controller struct {
	networksService networksService
	authMiddleware  *middlewares.AuthMiddleware
}

func New(
	networksService networksService,
	authMiddleware *middlewares.AuthMiddleware,
) *Controller {
	return &Controller{
		networksService: networksService,
		authMiddleware:  authMiddleware,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	netManagement := c.authMiddleware.AuthorizeNetManagement()
	router.Route("/networks", func(r chi.Router) {
		netManagementRoutes := r.With(netManagement)
		netManagementRoutes.Post("/", c.createNetwork)
	})
}
