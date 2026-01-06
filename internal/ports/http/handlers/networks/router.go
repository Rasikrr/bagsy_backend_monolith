package networks

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/command"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/go-chi/chi/v5"
)

type networksService interface {
	Create(ctx context.Context, createReq *command.CreateNetworkCommand) error
	GetByCode(ctx context.Context, code string) (*entity.Network, error)
}

type pointsService interface {
	GetByNetworkCode(ctx context.Context, networkCode string) ([]*entity.Point, error)
}
type Controller struct {
	networksService networksService
	pointsService   pointsService
	authMiddleware  *middlewares.AuthMiddleware
}

func New(
	networksService networksService,
	pointsService pointsService,
	authMiddleware *middlewares.AuthMiddleware,
) *Controller {
	return &Controller{
		networksService: networksService,
		pointsService:   pointsService,
		authMiddleware:  authMiddleware,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	netManagement := c.authMiddleware.AuthorizeNetManagement()
	router.Route("/api/v1/networks", func(r chi.Router) {
		netManagementRoutes := r.With(netManagement)
		netManagementRoutes.Post("/", c.createNetwork)

		r.Get("/{code}", c.getNetwork)
		r.Get("/{code}/points", c.getPointsByNetwork)
	})
}
