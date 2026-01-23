package servicecategories

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/go-chi/chi/v5"
)

type serviceCategoriesService interface {
	GetByPointCode(ctx context.Context, pointCode string) ([]*service.CategoryWithSubcategories, error)
}

type Controller struct {
	serviceCategoriesService serviceCategoriesService
}

func New(serviceCategoriesService serviceCategoriesService) *Controller {
	return &Controller{
		serviceCategoriesService: serviceCategoriesService,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	router.Route("/api/v1/service-categories", func(r chi.Router) {
		r.Get("/{point_code}", c.getByPointCode)
	})
}
