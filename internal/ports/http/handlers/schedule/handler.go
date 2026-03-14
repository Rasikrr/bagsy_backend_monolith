package schedule

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	domainSchedule "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/schedule"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type scheduleUseCase interface {
	GetLocationSchedule(ctx context.Context, orgCtx *access.OrgContext, locationID uuid.UUID, start, end time.Time) ([]*domainSchedule.LocationScheduleSlot, error)
	SetLocationSchedule(ctx context.Context, orgCtx *access.OrgContext, input uc.SetLocationScheduleInput) error
	DeleteLocationSchedule(ctx context.Context, orgCtx *access.OrgContext, locationID uuid.UUID, start, end time.Time) error
	GetEmployeeSchedule(ctx context.Context, orgCtx *access.OrgContext, employeeID uuid.UUID, start, end time.Time) ([]*domainSchedule.EmployeeScheduleSlot, error)
	SetEmployeeSchedule(ctx context.Context, orgCtx *access.OrgContext, input uc.SetEmployeeScheduleInput) error
	DeleteEmployeeSchedule(ctx context.Context, orgCtx *access.OrgContext, employeeID uuid.UUID, start, end time.Time) error
}

type Handler struct {
	scheduleUC    scheduleUseCase
	authMid       *middlewares.Auth
	orgContextMid *middlewares.OrgContext
}

func New(
	scheduleUC scheduleUseCase,
	authMid *middlewares.Auth,
	orgContextMid *middlewares.OrgContext,
) *Handler {
	return &Handler{
		scheduleUC:    scheduleUC,
		authMid:       authMid,
		orgContextMid: orgContextMid,
	}
}

func (h *Handler) Init(router *chi.Mux) {
	router.Route("/api/v1/location-schedules/{locationID}", func(r chi.Router) {
		r.Use(h.authMid.Handle)
		r.Use(h.orgContextMid.Handle)

		r.Get("/", h.getLocationSchedule)
		r.Put("/", h.setLocationSchedule)
		r.Delete("/", h.deleteLocationSchedule)
	})

	router.Route("/api/v1/employee-schedules/{employeeID}", func(r chi.Router) {
		r.Use(h.authMid.Handle)
		r.Use(h.orgContextMid.Handle)

		r.Get("/", h.getEmployeeSchedule)
		r.Put("/", h.setEmployeeSchedule)
		r.Delete("/", h.deleteEmployeeSchedule)
	})
}
