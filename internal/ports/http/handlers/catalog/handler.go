package catalog

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/catalog"
	"github.com/go-chi/chi/v5"
)

type catalogUseCase interface {
	CreateService(ctx context.Context, orgCtx *access.OrgContext, input uc.CreateServiceInput) (*uc.CreateServiceOutput, error)
	CreateEmployeeService(ctx context.Context, orgCtx *access.OrgContext, input uc.CreateEmployeeServiceInput) (*uc.CreateEmployeeServiceOutput, error)
}

type Handler struct {
	catalogUseCase catalogUseCase
	authMid        *middlewares.Auth
	orgContextMid  *middlewares.OrgContext
}

func New(
	catalogUC catalogUseCase,
	authMid *middlewares.Auth,
	orgContextMid *middlewares.OrgContext,
) *Handler {
	return &Handler{
		catalogUseCase: catalogUC,
		authMid:        authMid,
		orgContextMid:  orgContextMid,
	}
}

func (h *Handler) Init(router *chi.Mux) {
	router.Route("/api/v1/services", func(r chi.Router) {
		r.Use(h.authMid.Handle)
		r.Use(h.orgContextMid.Handle)

		r.Post("/", h.createService)
	})

	router.Route("/api/v1/employee-services", func(r chi.Router) {
		r.Use(h.authMid.Handle)
		r.Use(h.orgContextMid.Handle)

		r.Post("/", h.createEmployeeService)
	})
}
