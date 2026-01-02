package form

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/go-chi/chi/v5"
)

type formsService interface {
	Create(ctx context.Context, form *entity.Form) error
}

type Controller struct {
	formsService formsService
}

func New(formsService formsService) *Controller {
	return &Controller{
		formsService: formsService,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	router.Route("/api/v1/forms", func(r chi.Router) {
		r.Post("/", c.createClient)
	})
}
