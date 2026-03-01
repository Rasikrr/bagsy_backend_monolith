package employee

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	employeeUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/employee"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/invite"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type inviteUseCase interface {
	SendInvite(ctx context.Context, orgCtx *access.OrgContext, input uc.SendInviteInput) (*uc.SendInviteOutput, error)
	ConfirmInvite(ctx context.Context, input uc.ConfirmInviteInput) (*uc.TokensOutput, error)
	ResendInvite(ctx context.Context, orgCtx *access.OrgContext, input uc.ResendInviteInput) (*uc.ResendInviteOutput, error)
}

type employeeUseCase interface {
	GetProfile(ctx context.Context, employeeID uuid.UUID) (*employeeUC.ProfileOutput, error)
	GetList(ctx context.Context, orgCtx *access.OrgContext, filter *identity.EmployeeFilter) (*employeeUC.ListOutput, error)
}

type Handler struct {
	inviteUseCase inviteUseCase
	employeeUC    employeeUseCase
	authMid       *middlewares.Auth
	orgContextMid *middlewares.OrgContext
}

func New(
	inviteUC inviteUseCase,
	employeeUC employeeUseCase,
	authMid *middlewares.Auth,
	orgContextMid *middlewares.OrgContext,
) *Handler {
	return &Handler{
		inviteUseCase: inviteUC,
		employeeUC:    employeeUC,
		authMid:       authMid,
		orgContextMid: orgContextMid,
	}
}

func (h *Handler) Init(router *chi.Mux) {
	router.Route("/api/v1/employees", func(r chi.Router) {
		// Authenticated routes (require auth + orgContext)
		r.Group(func(r chi.Router) {
			r.Use(h.authMid.Handle)
			r.Use(h.orgContextMid.Handle)
			r.Get("/", h.getList)
			r.Get("/me", h.getMe)
			r.Post("/invite", h.sendInvite)
			r.Post("/invite/resend", h.resendInvite)
		})

		// Unauthenticated routes
		r.Post("/invite/confirm", h.confirmInvite)
	})
}
