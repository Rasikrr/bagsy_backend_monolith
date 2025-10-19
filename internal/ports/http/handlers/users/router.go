package users

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	usersService   users.Service
	authMiddleware middlewares.AuthMiddleware
}

func NewController(
	usersService users.Service,
	authMiddleware middlewares.AuthMiddleware) *Controller {
	return &Controller{
		usersService:   usersService,
		authMiddleware: authMiddleware,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	auth := c.authMiddleware.Handle

	router.Route("/api/v1/users", func(r chi.Router) {
		r.Get("/", auth(c.getByPhone))
		r.Put("/", auth(c.update))
	})

	minRoleMiddleware := c.authMiddleware.RequireMinRole(enum.RoleManager)

	router.Route("/api/v1/users/admin", func(r chi.Router) {
		r.Get("/phone/{phone}", minRoleMiddleware(c.getByPhoneAdmin))
		r.Get("/network/{network_code}", minRoleMiddleware(c.getByNetworkCode))
		r.Get("/point/{point_code}", minRoleMiddleware(c.getByPointCode))
	})
}
