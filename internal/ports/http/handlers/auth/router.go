package auth

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/command"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/go-chi/chi/v5"
)

type authService interface {
	Login(ctx context.Context, phone, password string) (access, refresh string, err error)
	Refresh(ctx context.Context, refreshToken string) (access, refresh string, err error)
	RegisterStaff(ctx context.Context, phone, pointCode string) error
	RegisterStaffConfirm(ctx context.Context, req *command.RegisterStaffConfirmRequest) (access, refresh string, err error)
}

type Controller struct {
	authService    authService
	authMiddleware middlewares.AuthMiddleware
}

func New(
	authService authService,
	authMiddleware middlewares.AuthMiddleware,
) *Controller {
	return &Controller{
		authService:    authService,
		authMiddleware: authMiddleware,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	managerMiddleware := c.authMiddleware.RequireRole(
		enum.RoleManager,
		enum.RoleNetManager,
		enum.RoleSelfOwner,
	)
	router.Route("/api/v1/auth", func(r chi.Router) {

		r.Post("/staff/register", managerMiddleware(c.registerStaff))
		r.Post("/staff/register/confirm", c.registerStaffConfirm)

		r.Post("/login", c.login)

		r.Post("/refresh", c.refresh)
	})
}
