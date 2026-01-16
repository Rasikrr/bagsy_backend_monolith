package networks

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/network"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/go-chi/chi/v5"
)

type networksService interface {
	GetByCode(ctx context.Context, code string) (*network.Network, error)
}

type pointsService interface {
	GetByNetworkCode(ctx context.Context, networkCode string) ([]*point.Point, error)
}
type Controller struct {
	networksService networksService
	pointsService   pointsService
	authMiddleware  *middlewares.Auth
}

func New(
	networksService networksService,
	pointsService pointsService,
	authMiddleware *middlewares.Auth,
) *Controller {
	return &Controller{
		networksService: networksService,
		pointsService:   pointsService,
		authMiddleware:  authMiddleware,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	router.Route("/api/v1/networks", func(r chi.Router) {
		r.Get("/{code}", c.getNetwork)
		r.Get("/{code}/points", c.getPointsByNetwork)
	})
}
