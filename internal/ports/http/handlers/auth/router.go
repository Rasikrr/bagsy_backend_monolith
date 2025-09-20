package auth

import (
	"github.com/Rasikrr/bugsy_backend_monolith/internal/services/auth"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/services/users"
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
		r.Post("/code", c.sendSmsCode)
		r.Post("/register", c.register)
		r.Post("/register/confirm", c.registerConfirm)
		r.Post("/login", c.login)
		r.Post("/refresh", c.refresh)
	})
}
