package auth

import (
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	// cache
	// tg
}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) Init(r *chi.Mux) {
	// group := r.Group("/api/v1")

	// TODO: resolve
	// group.GET("/users/sms", c.Handler(c.sendSmsCode))
}
