package booking

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/booking"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type bookingUseCase interface {
	GetAvailableSlots(ctx context.Context, input uc.GetAvailableSlotsInput) (*uc.GetAvailableSlotsOutput, error)
	Create(ctx context.Context, input uc.CreateBookingInput) (*uc.CreateBookingOutput, error)
	Confirm(ctx context.Context, appointmentID uuid.UUID, code string) error
	Cancel(ctx context.Context, appointmentID uuid.UUID, reason string) error
}

type Handler struct {
	bookingUC    bookingUseCase
	clientsMid   *middlewares.Clients
	employeesMid *middlewares.Employees
}

func New(
	bookingUC bookingUseCase,
	clientsMid *middlewares.Clients,
	employeesMid *middlewares.Employees,
) *Handler {
	return &Handler{
		bookingUC:    bookingUC,
		clientsMid:   clientsMid,
		employeesMid: employeesMid,
	}
}

func (h *Handler) Init(router *chi.Mux) {
	router.Route("/api/v1/bookings", func(r chi.Router) {
		// Публичные эндпоинты (или доступные клиентам)
		r.Get("/slots", h.getSlots)
		r.Post("/", h.create)
		r.Post("/{id}/confirm", h.confirm)

		// Эндпоинты для сотрудников организации
		r.Group(func(admin chi.Router) {
			admin.Use(h.employeesMid.Handle)
			admin.Post("/{id}/cancel", h.cancel)
		})
	})
}
