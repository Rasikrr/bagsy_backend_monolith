// nolint:revive
package swagger

import (
	"fmt"

	_ "github.com/Rasikrr/bugsy_backend_monolith/docs/swagger"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Controller struct {
	scheme string
	host   string
}

func New(scheme, host string) *Controller {
	return &Controller{
		scheme: scheme,
		host:   host,
	}
}

func (c *Controller) Init(router *chi.Mux) {
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s://%s/swagger/doc.json", c.scheme, c.host)),
	))
}
