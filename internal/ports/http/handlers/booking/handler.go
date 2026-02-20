package booking

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/booking"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type bookingUseCase interface {
	Create(ctx context.Context, input uc.CreateBookingInput) (*uc.CreateBookingOutput, error)
	Confirm(ctx context.Context, appointmentID uuid.UUID, code string) error
}

type Handler struct {
	bookingUC      bookingUseCase
	authMiddleware *middlewares.Auth
	orgContextMid  *middlewares.OrgContext
}

func New(
	bookingUC bookingUseCase,
	authMiddleware *middlewares.Auth,
	orgContextMid *middlewares.OrgContext,
) *Handler {
	return &Handler{
		bookingUC:      bookingUC,
		authMiddleware: authMiddleware,
		orgContextMid:  orgContextMid,
	}
}

func (h *Handler) Init(router *chi.Mux) {
	router.Route("/api/v1/bookings", func(r chi.Router) {
		// Публичные эндпоинты (или доступные клиентам)
		r.Post("/", h.create)
		r.Post("/{id}/confirm", h.confirm)
	})
}
