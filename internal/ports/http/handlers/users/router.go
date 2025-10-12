package users

import (
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
	router.Route("/api/v1/users", func(r chi.Router) {
		r.Post("/", c.authMiddleware.Handle(c.getByParams))
		r.Get("/{phone}", c.authMiddleware.Handle(c.getByPhone))
		r.Put("/{phone}", c.authMiddleware.Handle(c.update))
	})
}
