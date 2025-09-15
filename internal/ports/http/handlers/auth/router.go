package auth

import (
	"github.com/Rasikrr/bugsy_backend_monolith/internal/services/auth"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	authService auth.Service
}

func New(authService auth.Service) *Controller {
	return &Controller{
		authService: authService,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	router.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/code", c.sendSmsCode)
	})
}
