package form

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/forms"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	formsService forms.Service
}

func NewController(formsService forms.Service) *Controller {
	return &Controller{
		formsService: formsService,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	router.Route("/api/v1/forms", func(r chi.Router) {
		r.Post("/", c.createClient)
	})
}
