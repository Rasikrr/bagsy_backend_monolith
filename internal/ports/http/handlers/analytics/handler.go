package analytics

import (
	"context"
	"net/http"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	domainAnalytics "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/analytics"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type analyticsUseCase interface {
	GetOverview(ctx context.Context, orgCtx *access.OrgContext, from, to time.Time, locationID *uuid.UUID) (*domainAnalytics.OverviewReport, error)
	GetMe(ctx context.Context, orgCtx *access.OrgContext, from, to time.Time) (*domainAnalytics.MeReport, error)
	GetStaff(ctx context.Context, orgCtx *access.OrgContext, from, to time.Time, locationID *uuid.UUID) (*domainAnalytics.StaffReport, error)
	GetStaffDetail(ctx context.Context, orgCtx *access.OrgContext, employeeID uuid.UUID, from, to time.Time) (*domainAnalytics.StaffDetailReport, error)
	GetLocation(ctx context.Context, orgCtx *access.OrgContext, locationID uuid.UUID, from, to time.Time) (*domainAnalytics.OverviewReport, error)
	GetFinance(ctx context.Context, orgCtx *access.OrgContext, from, to time.Time, locationID *uuid.UUID) (*domainAnalytics.FinanceReport, error)
	GetClients(ctx context.Context, orgCtx *access.OrgContext, from, to time.Time, locationID *uuid.UUID) (*domainAnalytics.ClientsReport, error)
}

type Handler struct {
	analyticsUC   analyticsUseCase
	authMid       *middlewares.Auth
	orgContextMid *middlewares.OrgContext
}

func New(
	analyticsUC analyticsUseCase,
	authMid *middlewares.Auth,
	orgContextMid *middlewares.OrgContext,
) *Handler {
	return &Handler{
		analyticsUC:   analyticsUC,
		authMid:       authMid,
		orgContextMid: orgContextMid,
	}
}

func (h *Handler) Init(router *chi.Mux) {
	router.Route("/api/v1/analytics", func(r chi.Router) {
		r.Use(h.authMid.Handle)
		r.Use(h.orgContextMid.Handle)

		r.Get("/overview", h.getOverview)
		r.Get("/me", h.getMe)
		r.Get("/staff", h.getStaff)
		r.Get("/staff/{employeeID}", h.getStaffDetail)
		r.Get("/locations/{locationID}", h.getLocation)
		r.Get("/finance", h.getFinance)
		r.Get("/clients", h.getClients)
	})
}

// setCacheHeaders проставляет приватное кэширование (фронт держит staleTime=5min).
func setCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "private, max-age=60")
}
