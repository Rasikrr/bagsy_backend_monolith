package calendar

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/bagsies"
	"github.com/go-chi/chi/v5"
)

type calendarService interface {
	GetCalendar(
		ctx context.Context,
		query *bagsies.GetCalendarQuery,
	) ([]*bagsies.CalendarElement, error)
}

type Controller struct {
	calendarService calendarService
	authMiddleware  *middlewares.Auth
}

func New(
	calendarService calendarService,
	authMiddleware *middlewares.Auth,
) *Controller {
	return &Controller{
		calendarService: calendarService,
		authMiddleware:  authMiddleware,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	auth := c.authMiddleware.Handle

	router.Route("/api/v1/calendar", func(r chi.Router) {
		r.With(auth).
			Get("/", c.getCalendar)
	})
}
