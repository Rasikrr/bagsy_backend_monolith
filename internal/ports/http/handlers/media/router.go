package media

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/dto"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/go-chi/chi/v5"
)

type mediaService interface {
	GenerateUploadURL(ctx context.Context, key string, contentType, purpose string) (*dto.UploadMediaResponse, error)
}

type Controller struct {
	mediaService   mediaService
	authMiddleware *middlewares.Auth
}

func New(
	mediaService mediaService,
	authMiddleware *middlewares.Auth,
) *Controller {
	return &Controller{
		mediaService:   mediaService,
		authMiddleware: authMiddleware,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	auth := c.authMiddleware.Handle
	router.Route("/api/v1/media", func(r chi.Router) {
		r.With(auth).
			Post("/upload", c.getUploadURL)
	})
}
