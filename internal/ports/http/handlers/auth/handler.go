package auth

import (
	"context"

	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/auth"
	"github.com/go-chi/chi/v5"
)

type registerOwnerUseCase interface {
	Register(ctx context.Context, req uc.RegisterInput) (*uc.RegisterOutput, error)
	Resend(ctx context.Context, req uc.ResendInput) (*uc.ResendOutput, error)
}

// Handler serves owner registration HTTP endpoints.
type Handler struct {
	registerOwnerUseCase registerOwnerUseCase
}

func New(
	registerOwnerUseCase registerOwnerUseCase,
) *Handler {
	return &Handler{
		registerOwnerUseCase: registerOwnerUseCase,
	}
}

func (h *Handler) Init(router *chi.Mux) {

}
