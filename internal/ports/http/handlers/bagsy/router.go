package bagsy

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/command"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/go-chi/chi/v5"
)

type bagsiesService interface {
	Create(ctx context.Context, req *command.CreateBagsyCommand) error
}

type Controller struct {
	bagsiesService bagsiesService
	authMiddleware middlewares.AuthMiddleware
}

func New(
	bagsiesService bagsiesService,
	authMiddleware middlewares.AuthMiddleware,
) *Controller {
	return &Controller{
		bagsiesService: bagsiesService,
		authMiddleware: authMiddleware,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	router.Route("/api/v1/bagsies", func(r chi.Router) {
		r.Post("/", c.createBagsy)
	})
}
