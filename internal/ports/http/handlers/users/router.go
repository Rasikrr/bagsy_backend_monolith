package users

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	usersService users.Service
}

func NewController(usersService users.Service) *Controller {
	return &Controller{
		usersService: usersService,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	router.Route("/api/v1/users", func(r chi.Router) {
		r.Get("/{phone}", c.getByPhone)
		r.Put("/{phone}", c.update)
	})
}
