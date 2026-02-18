package auth

import (
	"context"

	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/auth"
	"github.com/go-chi/chi/v5"
)

type registerOwnerUseCase interface {
	Register(ctx context.Context, req uc.RegisterInput) (*uc.RegisterOutput, error)
	Resend(ctx context.Context, req uc.ResendInput) (*uc.ResendOutput, error)
	VerifyRegistration(ctx context.Context, req uc.VerifyInput) (*uc.TokensOutput, error)
}

type authUseCase interface {
	LoginEmployee(ctx context.Context, phone, password string) (*uc.TokensOutput, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*uc.TokensOutput, error)
	Logout(ctx context.Context, token string) error
}

type resetPasswordUseCase interface {
	RequestReset(ctx context.Context, req uc.RequestResetInput) error
	ConfirmReset(ctx context.Context, req uc.ConfirmResetInput) error
}

// Handler serves auth HTTP endpoints.
type Handler struct {
	registerOwnerUseCase registerOwnerUseCase
	authUseCase          authUseCase
	resetPasswordUseCase resetPasswordUseCase
}

func New(
	registerOwnerUseCase registerOwnerUseCase,
	authUseCase authUseCase,
	resetPasswordUseCase resetPasswordUseCase,
) *Handler {
	return &Handler{
		registerOwnerUseCase: registerOwnerUseCase,
		authUseCase:          authUseCase,
		resetPasswordUseCase: resetPasswordUseCase,
	}
}

func (h *Handler) Init(router *chi.Mux) {
	router.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", h.registerOwner)
		r.Post("/register/verify", h.verifyOwner)
		r.Post("/register/resend", h.resendOwner)

		r.Post("/login", h.loginEmployee)
		r.Post("/refresh", h.refreshTokens)
		r.Post("/logout", h.logout)

		r.Post("/password/reset", h.requestPasswordReset)
		r.Post("/password/reset/confirm", h.confirmPasswordReset)
	})
}
