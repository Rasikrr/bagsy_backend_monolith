package points

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/go-chi/chi/v5"
)

type pointsService interface {
	Create(ctx context.Context, point *entity.Point) error
	GetByCode(ctx context.Context, code string) (*entity.Point, error)
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
		// TODO:  по идее только нет менеджер/ самозанятые могут создавать точки
		netManagementsRoutes := r.With(netManagement)
		netManagementsRoutes.Post("/", c.createPoint)
		// TODO: Вроде все челы могут получать?
		r.Get("/{code}", c.getPoint)
	})
}
