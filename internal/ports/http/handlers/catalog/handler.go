package catalog

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	catalogDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/catalog"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type catalogUseCase interface {
	CreateService(ctx context.Context, orgCtx *access.OrgContext, input uc.CreateServiceInput) (*uc.CreateServiceOutput, error)
	UpdateService(ctx context.Context, orgCtx *access.OrgContext, input uc.UpdateServiceInput) error
	DeleteService(ctx context.Context, orgCtx *access.OrgContext, serviceID uuid.UUID) error
	CreateEmployeeService(ctx context.Context, orgCtx *access.OrgContext, input uc.CreateEmployeeServiceInput) (*uc.CreateEmployeeServiceOutput, error)
	DeleteEmployeeService(ctx context.Context, orgCtx *access.OrgContext, employeeServiceID uuid.UUID) error
	GetServiceCategories(ctx context.Context, locationCategoryID uuid.UUID) ([]uc.ServiceCategoryTree, error)
	GetServicesByLocation(ctx context.Context, locationID uuid.UUID) ([]*catalogDomain.Service, error)
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
		r.Get("/{id}", h.getServicesByLocation)

		r.Group(func(r chi.Router) {
			r.Use(h.authMid.Handle)
			r.Use(h.orgContextMid.Handle)

			r.Post("/", h.createService)
			r.Put("/{id}", h.updateService)
			r.Delete("/{id}", h.deleteService)
		})
	})

	router.Route("/api/v1/employee-services", func(r chi.Router) {
		r.Use(h.authMid.Handle)
		r.Use(h.orgContextMid.Handle)

		r.Post("/", h.createEmployeeService)
		r.Delete("/{id}", h.deleteEmployeeService)
	})

	router.Route("/api/v1/service-categories", func(r chi.Router) {
		r.Use(h.authMid.Handle)
		r.Use(h.orgContextMid.Handle)

		r.Get("/", h.getServiceCategories)
	})
}
