package pointcategories

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/go-chi/chi/v5"
)

type pointCategoriesService interface {
	GetAll(ctx context.Context) ([]*point.Category, error)
}

type Controller struct {
	pointCategoriesService pointCategoriesService
}

func New(pointCategoriesService pointCategoriesService) *Controller {
	return &Controller{
		pointCategoriesService: pointCategoriesService,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	router.Route("/api/v1/point-categories", func(r chi.Router) {
		r.Get("/", c.getCategories)
	})
}
