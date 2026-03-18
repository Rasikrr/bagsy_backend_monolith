package booking

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/booking"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/booking"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type bookingUseCase interface {
	Create(ctx context.Context, input uc.CreateBookingInput) (*uc.CreateBookingOutput, error)
	CreateDirect(ctx context.Context, orgCtx *access.OrgContext, input uc.CreateBookingInput) (*uc.CreateBookingOutput, error)
	Confirm(ctx context.Context, appointmentID uuid.UUID, code string) error
	Cancel(ctx context.Context, orgCtx *access.OrgContext, appointmentID uuid.UUID, reason string) error
	ResendOTP(ctx context.Context, appointmentID uuid.UUID) error
	GetAvailableSlots(ctx context.Context, input uc.GetAvailableSlotsInput) (*uc.GetAvailableSlotsOutput, error)
	GetCalendar(ctx context.Context, orgCtx *access.OrgContext, input uc.GetCalendarInput) ([]booking.CalendarEntry, error)
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
	router.Route("/api/v1/appointments", func(r chi.Router) {
		// Публичные эндпоинты (или доступные клиентам)
		r.Post("/", h.create)
		r.Post("/slots", h.getSlots)
		r.Post("/{id}/confirm", h.confirm)
		r.Post("/{id}/resend-otp", h.resendOTP)

		// Эндпоинты для сотрудников организации
		r.Group(func(admin chi.Router) {
			admin.Use(
				h.authMiddleware.Handle,
				h.orgContextMid.Handle,
			)

			admin.Post("/direct", h.createDirect)
			admin.Post("/{id}/cancel", h.cancel)
			admin.Get("/calendar", h.getCalendar)
		})
	})
}
