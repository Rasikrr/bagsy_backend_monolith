package employee

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/invite"
	"github.com/go-chi/chi/v5"
)

type inviteUseCase interface {
	SendInvite(ctx context.Context, orgCtx *access.OrgContext, input uc.SendInviteInput) (*uc.SendInviteOutput, error)
	ConfirmInvite(ctx context.Context, input uc.ConfirmInviteInput) (*uc.TokensOutput, error)
	ResendInvite(ctx context.Context, orgCtx *access.OrgContext, input uc.ResendInviteInput) (*uc.ResendInviteOutput, error)
	VerifyInviteToken(ctx context.Context, token string) (*uc.VerifyInviteTokenOutput, error)
}

type Handler struct {
	inviteUseCase inviteUseCase
	authMid       *middlewares.Auth
	orgContextMid *middlewares.OrgContext
}

func New(
	inviteUC inviteUseCase,
	authMid *middlewares.Auth,
	orgContextMid *middlewares.OrgContext,
) *Handler {
	return &Handler{
		inviteUseCase: inviteUC,
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
			r.Post("/invite", h.sendInvite)
			r.Post("/invite/resend", h.resendInvite)
		})

		// Unauthenticated routes
		r.Post("/invite/confirm", h.confirmInvite)
		r.Get("/invite/verify/{token}", h.verifyInviteToken)
	})
}
