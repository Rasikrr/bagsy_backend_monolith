package location

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/location"
	"github.com/go-chi/chi/v5"
)

type locationUseCase interface {
	Create(ctx context.Context, orgCtx *access.OrgContext, input uc.CreateLocationInput) (*uc.CreateLocationOutput, error)
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
		r.Use(h.authMid.Handle)
		r.Use(h.orgContextMid.Handle)

		r.Post("/", h.create)
	})
}
