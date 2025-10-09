package bagsies

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/bagsies"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	bagsyService bagsies.Service
	authService  auth.Service
	usersService users.Service
	authMW       middlewares.AuthMiddleware
}

func New(
	bagsyService bagsies.Service,
	authService auth.Service,
	usersService users.Service,
	authMiddleware middlewares.AuthMiddleware,
) *Controller {
	return &Controller{
		bagsyService: bagsyService,
		authService:  authService,
		usersService: usersService,
		authMW:       authMiddleware,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	router.Route("/api/v1/bagsies", func(r chi.Router) {
		r.Post("/create", c.create)
	})
}
