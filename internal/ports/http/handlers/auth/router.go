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
	Logout(ctx context.Context, refreshToken string) error
	RegisterManagement(ctx context.Context, req *command.RegisterManagementCommand) error
	RegisterManagementConfirm(ctx context.Context, phone, code string) (access, refresh string, err error)
	ResendRegisterManagementCode(ctx context.Context, phone string) error
	RegisterStaff(ctx context.Context, req *command.RegisterStaffCommand) error
	RegisterStaffConfirm(ctx context.Context, req *command.RegisterStaffConfirmCommand) (access, refresh string, err error)
	SendPasswordChangeLink(ctx context.Context, phone string) error
	ChangePassword(ctx context.Context, req *command.ChangePasswordConfirmCommand) error
}

type Controller struct {
	authService    authService
	authMiddleware *middlewares.AuthMiddleware
}

func New(
	authService authService,
	authMiddleware *middlewares.AuthMiddleware,
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
		r.Route("/staff", func(r chi.Router) {
			r.With(managerMiddleware).
				Post("/register", c.registerStaff)

			r.Post("/register/confirm", c.registerStaffConfirm)
		})

		r.Route("/management", func(r chi.Router) {
			r.Post("/register", c.registerManagement)
			r.Post("/register/resend", c.resendRegisterManagement)
			r.Post("/register/confirm", c.registerManagementConfirm)
		})

		r.Post("/password/change/request", c.changePassword)
		r.Post("/password/change/confirm", c.changePasswordConfirm)
		r.Post("/login", c.login)
		r.Post("/refresh", c.refresh)
		r.Post("/logout", c.logout)
	})
}
