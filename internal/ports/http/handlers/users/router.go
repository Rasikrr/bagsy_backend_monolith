package users

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/query"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/go-chi/chi/v5"
)

type userService interface {
	GetUserProfile(ctx context.Context) (*user.User, error)
	GetListByFilter(ctx context.Context, filter *user.Filter) (*query.Page[*user.User], error)
	UpdateProfile(ctx context.Context, cmd *user.UpdateUserCommand) (*user.User, error)
	UpdateSchedule(ctx context.Context, phone string, schedule user.Schedule) error
	RemoveAvatar(ctx context.Context) error
}

type Controller struct {
	userService    userService
	authMiddleware *middlewares.Auth
}

func New(
	userService userService,
	authMiddleware *middlewares.Auth,
) *Controller {
	return &Controller{
		userService:    userService,
		authMiddleware: authMiddleware,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	auth := c.authMiddleware.Handle
	management := c.authMiddleware.AuthorizeManagement()

	router.Route("/api/v1/staff", func(r chi.Router) {
		managersRoutes := r.With(management)
		managersRoutes.Get("/", c.getUsers)
	})

	router.Route("/api/v1/users", func(r chi.Router) {
		authenticated := r.With(auth)
		authenticated.Get("/me", c.getMyProfile)
		authenticated.Put("/me", c.updateUser)
		authenticated.Put("/me/schedule", c.updateSchedule)
		authenticated.Delete("/me/avatar", c.deleteAvatar)
	})
}
