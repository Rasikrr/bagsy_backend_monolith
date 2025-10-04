package auth

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	authService  auth.Service
	usersService users.Service
}

func New(
	authService auth.Service,
	usersService users.Service,
) *Controller {
	return &Controller{
		authService:  authService,
		usersService: usersService,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	router.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", c.register)
		r.Post("/register/confirm", c.registerConfirm)
		r.Post("/login", c.login)
		r.Post("/logout", c.logout)
		r.Post("/refresh", c.refresh)
	})
}
