package media

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/media"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type mediaUseCase interface {
	GenerateUploadURL(ctx context.Context, input uc.GenerateUploadURLInput) (*uc.GenerateUploadURLOutput, error)
	ConfirmUpload(ctx context.Context, assetID uuid.UUID) error
}

type Handler struct {
	mediaUseCase mediaUseCase
	authMid      *middlewares.Auth
}

func New(
	mediaUC mediaUseCase,
	authMid *middlewares.Auth,
) *Handler {
	return &Handler{
		mediaUseCase: mediaUC,
		authMid:      authMid,
	}
}

func (h *Handler) Init(router *chi.Mux) {
	router.Route("/api/v1/media", func(r chi.Router) {
		r.Use(h.authMid.Handle)

		r.Post("/upload", h.upload)
		r.Post("/{id}/confirm", h.confirm)
	})
}
