package bagsies

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/command"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type bagsiesService interface {
	Create(ctx context.Context, req *command.CreateBagsyCommand) (uuid.UUID, error)
	Confirm(ctx context.Context, bagsyID uuid.UUID, code string) error
	ResendConfirmationCode(ctx context.Context, bagsyID uuid.UUID) error
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
		r.Post("/resent", c.resendCode)
		r.Post("/confirm", c.confirmBagsy)
	})
}
