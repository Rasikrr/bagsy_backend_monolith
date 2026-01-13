package users

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/command"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/dto"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/query"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/go-chi/chi/v5"
)

type userService interface {
	GetUserProfile(ctx context.Context) (*dto.UserWithAvatar, error)
	GetListByFilter(ctx context.Context, filter *query.UserFilter) (*dto.PaginatedUsers, error)
	UpdateProfile(ctx context.Context, cmd *command.UpdateUserCommand) (*dto.UserWithAvatar, error)
	UpdateSchedule(ctx context.Context, phone string, schedule []entity.StaffSchedule) error
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
