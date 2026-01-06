package users

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/query"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/go-chi/chi/v5"
)

type userService interface {
	GetByFilter(ctx context.Context, filter *query.UserFilter) ([]*entity.User, error)
}

type Controller struct {
	userService    userService
	authMiddleware *middlewares.AuthMiddleware
}

func New(
	userService userService,
	authMiddleware *middlewares.AuthMiddleware,
) *Controller {
	return &Controller{
		userService:    userService,
		authMiddleware: authMiddleware,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	management := c.authMiddleware.AuthorizeManagement()
	router.Route("/api/v1/staff", func(r chi.Router) {
		managersRoutes := r.With(management)
		managersRoutes.Get("/", c.getUsers)
	})
}
