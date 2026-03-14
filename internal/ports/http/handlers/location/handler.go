package location

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	domainLoc "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/location"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type locationUseCase interface {
	Create(ctx context.Context, orgCtx *access.OrgContext, input uc.CreateLocationInput) (*uc.CreateLocationOutput, error)
	GetList(ctx context.Context, orgCtx *access.OrgContext, filter *domainLoc.Filter) (*shared.Page[*domainLoc.Location], error)
	GetByID(ctx context.Context, orgCtx *access.OrgContext, id uuid.UUID) (*domainLoc.Location, error)
	GetCategories(ctx context.Context) ([]*domainLoc.Category, error)
	UpdateLocation(ctx context.Context, orgCtx *access.OrgContext, input uc.UpdateLocationInput) error
	DeleteLocation(ctx context.Context, orgCtx *access.OrgContext, locationID uuid.UUID) error
}

type Handler struct {
	locationUseCase locationUseCase
	authMid         *middlewares.Auth
	orgContextMid   *middlewares.OrgContext
}

func New(
	createUC locationUseCase,
	authMid *middlewares.Auth,
	orgContextMid *middlewares.OrgContext,
) *Handler {
	return &Handler{
		locationUseCase: createUC,
		authMid:         authMid,
		orgContextMid:   orgContextMid,
	}
}

func (h *Handler) Init(router *chi.Mux) {
	router.Route("/api/v1/locations", func(r chi.Router) {
		r.Get("/categories", h.getCategories)

		r.Group(func(r chi.Router) {
			r.Use(h.authMid.Handle)
			r.Use(h.orgContextMid.Handle)

			r.Get("/", h.getList)
			r.Get("/{id}", h.getByID)
			r.Post("/", h.create)
			r.Put("/{id}", h.updateLocation)
			r.Delete("/{id}", h.deleteLocation)
		})
	})
}
