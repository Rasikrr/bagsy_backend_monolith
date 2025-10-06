package auth

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	authService    auth.Service
	usersService   users.Service
	authMiddleware middlewares.AuthMiddleware
}

func New(
	authService auth.Service,
	usersService users.Service,
	authMiddleware middlewares.AuthMiddleware,
) *Controller {
	return &Controller{
		authService:    authService,
		usersService:   usersService,
		authMiddleware: authMiddleware,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	router.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", c.authMiddleware.Handle(c.register))
		r.Post("/register/confirm", c.registerConfirm)
		r.Post("/login", c.login)
		r.Post("/refresh", c.refresh)
	})
}
