package auth

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	authS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"
	"github.com/go-chi/chi/v5"
)

const (
	rateLimiterLimit       = 3
	rateLimiterPerDuration = time.Minute
)

type authService interface {
	Login(ctx context.Context, phone, password string) (access, refresh string, err error)
	Refresh(ctx context.Context, refreshToken string) (access, refresh string, err error)
	Logout(ctx context.Context, refreshToken string) error
	InspectAuthToken(ctx context.Context, token string) (*authS.InviteTokenInfo, error)
	RegisterManagement(ctx context.Context, req *auth.RegisterManagementCommand) error
	RegisterManagementConfirm(ctx context.Context, phone, code string) (access, refresh string, err error)
	ResendRegisterManagementCode(ctx context.Context, phone string) error
	RegisterStaff(ctx context.Context, req *auth.RegisterStaffCommand) error
	RegisterStaffConfirm(ctx context.Context, req *auth.RegisterStaffConfirmCommand) (access, refresh string, err error)
	ResendRegisterStaffLink(ctx context.Context, phone string) error
	SendPasswordChangeLink(ctx context.Context, phone string) error
	ChangePassword(ctx context.Context, req *auth.ChangePasswordConfirmCommand) error
}

type Controller struct {
	authService        authService
	authMiddleware     *middlewares.Auth
	rateLimiterFactory *middlewares.RateLimiterFactory
}

func New(
	authService authService,
	authMiddleware *middlewares.Auth,
	rateLimiterFactory *middlewares.RateLimiterFactory,
) *Controller {
	return &Controller{
		authService:        authService,
		authMiddleware:     authMiddleware,
		rateLimiterFactory: rateLimiterFactory,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	managerMiddleware := c.authMiddleware.RequireRole(
		user.RoleManager,
		user.RoleNetManager,
		user.RoleSelfOwner,
	)

	rateLimiter := c.rateLimiterFactory.NewRateLimiter(rateLimiterLimit, rateLimiterPerDuration)

	router.Route("/api/v1/auth", func(r chi.Router) {
		r.Route("/staff", func(r chi.Router) {
			r.With(managerMiddleware).
				Post("/register", c.registerStaff)
			r.With(managerMiddleware, rateLimiter).
				Post("/register/resend", c.registerStaffResend)

			r.Post("/register/confirm", c.registerStaffConfirm)
		})

		r.Route("/management", func(r chi.Router) {
			r.Post("/register", c.registerManagement)

			r.With(rateLimiter).
				Post("/register/resend", c.resendRegisterManagement)
			r.Post("/register/confirm", c.registerManagementConfirm)
		})

		r.With(rateLimiter).
			Post("/password/change", c.changePassword)
		r.Post("/password/change/confirm", c.changePasswordConfirm)

		r.Get("/verify-auth-token/{token}", c.verifyAuthToken)
		r.Post("/login", c.login)
		r.Post("/refresh", c.refresh)
		r.Post("/logout", c.logout)
	})
}
