// nolint
package points

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/go-chi/chi/v5"
)

type pointsService interface {
	Create(ctx context.Context, cmd *point.CreatePointCommand) (*point.Point, error)
	GetPublicByCode(ctx context.Context, code string) (*point.Point, error)
}

type Controller struct {
	pointsService  pointsService
	authMiddleware *middlewares.Auth
}

func New(
	pointsService pointsService,
	authMiddleware *middlewares.Auth,
) *Controller {
	return &Controller{
		pointsService:  pointsService,
		authMiddleware: authMiddleware,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	netManagement := c.authMiddleware.AuthorizeNetManagement()
	router.Route("/api/v1/points", func(r chi.Router) {
		netManagementsRoutes := r.With(netManagement)

		netManagementsRoutes.Post("/", c.createPoint)

		r.Get("/{code}", c.getPoint)
	})
}
