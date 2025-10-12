package points

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/points"
	"github.com/go-chi/chi/v5"
)

// Controller для points
// Можно расширить при необходимости

type Controller struct {
	service points.Service
}

func NewController(service points.Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) Init(router *chi.Mux) {
	router.Route("/api/v1/points", func(r chi.Router) {
		r.Get("/{code}", c.getByCode)
		r.Put("/{code}", c.updateByCode)
		r.Post("/", c.create)
		r.Delete("/{code}", c.deleteByCode)
	})
}
