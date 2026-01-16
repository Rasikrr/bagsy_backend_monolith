package bagsies

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/bagsy"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type bagsiesService interface {
	Create(ctx context.Context, req *bagsy.CreateBagsyCommand) (uuid.UUID, error)
	Confirm(ctx context.Context, bagsyID uuid.UUID, code string) error
	ResendConfirmationCode(ctx context.Context, bagsyID uuid.UUID) error
	GetAvailableSlots(ctx context.Context, cmd *bagsy.GetAvailableSlotsCommand) (*bagsy.AvailableSlots, error)
}

type Controller struct {
	bagsiesService bagsiesService
	authMiddleware *middlewares.Auth
}

func New(
	bagsiesService bagsiesService,
	authMiddleware *middlewares.Auth,
) *Controller {
	return &Controller{
		bagsiesService: bagsiesService,
		authMiddleware: authMiddleware,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	router.Route("/api/v1/bagsies", func(r chi.Router) {
		r.Post("/", c.createBagsy)
		r.Post("/resend", c.resendCode)
		r.Post("/confirm", c.confirmBagsy)
		r.Post("/slots", c.getSlots)
		r.Post("/slots/day", c.getSlotsForDay)
	})
}
