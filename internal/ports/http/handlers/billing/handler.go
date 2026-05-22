package billing

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	domainBilling "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/billing"
	"github.com/go-chi/chi/v5"
)

type billingUseCase interface {
	GetSubscription(ctx context.Context, orgCtx *access.OrgContext) (*uc.SubscriptionOutput, error)
	Activate(ctx context.Context, orgCtx *access.OrgContext, input uc.ActivateInput) error
	RequestCancellation(ctx context.Context, orgCtx *access.OrgContext) error
	UndoCancellation(ctx context.Context, orgCtx *access.OrgContext) error
	ListPlans(ctx context.Context) ([]*domainBilling.Plan, error)
}

type Handler struct {
	billingUC     billingUseCase
	authMid       *middlewares.Auth
	orgContextMid *middlewares.OrgContext
}

func New(
	billingUC billingUseCase,
	authMid *middlewares.Auth,
	orgContextMid *middlewares.OrgContext,
) *Handler {
	return &Handler{
		billingUC:     billingUC,
		authMid:       authMid,
		orgContextMid: orgContextMid,
	}
}

func (h *Handler) Init(router *chi.Mux) {
	router.Route("/api/v1/plans", func(r chi.Router) {
		r.Get("/", h.listPlans)
	})

	router.Route("/api/v1/subscription", func(r chi.Router) {
		r.Use(h.authMid.Handle)
		r.Use(h.orgContextMid.Handle)

		r.Get("/", h.getSubscription)
		r.Post("/activate", h.activate)
		r.Post("/cancel", h.requestCancellation)
		r.Post("/undo-cancel", h.undoCancellation)
	})
}
