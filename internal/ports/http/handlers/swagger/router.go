// nolint:revive
package swagger

import (
	_ "github.com/Rasikrr/bugsy_backend_monolith/docs/swagger"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Controller struct {
}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) Init(router *chi.Mux) {
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
}
