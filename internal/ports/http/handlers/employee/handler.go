package employee

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	catalogDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	employeeUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/employee"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/invite"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// extractOrgContextAndID is a helper that extracts OrgContext from the request context
// and parses the employee UUID from the "id" URL parameter.
// Returns false if either extraction fails (response is already written).

type inviteUseCase interface {
	SendInvite(ctx context.Context, orgCtx *access.OrgContext, input uc.SendInviteInput) (*uc.SendInviteOutput, error)
	ConfirmInvite(ctx context.Context, input uc.ConfirmInviteInput) (*uc.TokensOutput, error)
	ResendInvite(ctx context.Context, orgCtx *access.OrgContext, input uc.ResendInviteInput) (*uc.ResendInviteOutput, error)
}

type employeeUseCase interface {
	GetProfile(ctx context.Context, employeeID uuid.UUID) (*employeeUC.ProfileOutput, error)
	GetList(ctx context.Context, orgCtx *access.OrgContext, filter *identity.EmployeeFilter) (*employeeUC.ListOutput, error)
	UpdateProfile(ctx context.Context, employeeID uuid.UUID, input employeeUC.UpdateProfileInput) (*employeeUC.ProfileOutput, error)
	TransferEmployee(ctx context.Context, orgCtx *access.OrgContext, employeeID uuid.UUID, input employeeUC.TransferInput) error
	ActivateEmployee(ctx context.Context, orgCtx *access.OrgContext, employeeID uuid.UUID) error
	DeactivateEmployee(ctx context.Context, orgCtx *access.OrgContext, employeeID uuid.UUID) error
	ChangeRole(ctx context.Context, orgCtx *access.OrgContext, employeeID uuid.UUID, input employeeUC.ChangeRoleInput) error
	ChangePermissions(ctx context.Context, orgCtx *access.OrgContext, employeeID uuid.UUID, input employeeUC.ChangePermissionsInput) error
}

type catalogUseCase interface {
	GetServicesByEmployee(ctx context.Context, employeeID uuid.UUID) ([]*catalogDomain.Service, error)
}

type Handler struct {
	inviteUseCase  inviteUseCase
	employeeUC     employeeUseCase
	catalogUseCase catalogUseCase
	authMid        *middlewares.Auth
	orgContextMid  *middlewares.OrgContext
}

func New(
	inviteUC inviteUseCase,
	employeeUC employeeUseCase,
	catalogUC catalogUseCase,
	authMid *middlewares.Auth,
	orgContextMid *middlewares.OrgContext,
) *Handler {
	return &Handler{
		inviteUseCase:  inviteUC,
		employeeUC:     employeeUC,
		catalogUseCase: catalogUC,
		authMid:        authMid,
		orgContextMid:  orgContextMid,
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
			r.Put("/me", h.updateMe)
			r.Post("/invite", h.sendInvite)
			r.Post("/invite/resend", h.resendInvite)

			r.Post("/{id}/transfer", h.transferEmployee)
			r.Post("/{id}/activate", h.activateEmployee)
			r.Post("/{id}/deactivate", h.deactivateEmployee)
			r.Patch("/{id}/role", h.changeRole)
			r.Patch("/{id}/permissions", h.changePermissions)
		})

		// Unauthenticated routes
		r.Post("/invite/confirm", h.confirmInvite)
		r.Get("/{id}/services", h.getEmployeeServices)
	})
}
