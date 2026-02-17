// nolint:revive
package swagger

import (
	"fmt"

	_ "github.com/Rasikrr/bagsy_backend_monolith/docs/swagger"
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
		httpSwagger.BeforeScript(`
			window.onload = function() {
				// Переопределяем метод авторизации для автоматической подстановки Bearer
				const originalAuthorize = window.ui.authActions.authorize;
				window.ui.authActions.authorize = function(payload) {
					if (payload.Bearer && payload.Bearer.value) {
						// Добавляем Bearer префикс если его нет
						const token = payload.Bearer.value.trim();
						if (!token.toLowerCase().startsWith('bearer ')) {
							payload.Bearer.value = 'Bearer ' + token;
						}
					}
					return originalAuthorize.call(this, payload);
				};
			}
		`),
	))
}
