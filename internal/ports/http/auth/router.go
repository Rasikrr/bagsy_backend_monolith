package auth

import (
	fasthttprouter "github.com/fasthttp/router"
)

type Controller struct {
	// cache
	// tg
}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) Init(r *fasthttprouter.Router) {
	// group := r.Group("/api/v1")

	// TODO: resolve
	// group.GET("/users/sms", c.Handler(c.sendSmsCode))
}
