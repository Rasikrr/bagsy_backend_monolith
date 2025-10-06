// nolint
package http

import (
	"context"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"net/http"

	docs "github.com/Rasikrr/bagsy_backend_monolith/docs/swagger"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/bagsies"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/form"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/swagger"
	authS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"
	bagsiesS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/bagsies"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/forms"
	usersS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
	"github.com/Rasikrr/core/enum"
	coreHttp "github.com/Rasikrr/core/http"
	"github.com/Rasikrr/core/log"
	"github.com/Rasikrr/core/version"
)

// @title Bagsy API
// @version 1.0
// @description API для приложения Bagsy
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@bagsy.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewServer(
	server *coreHttp.Server,
	swaggerHost, swaggerScheme string,
	authService authS.Service,
	formsService forms.Service,
	usersService usersS.Service,
	bagsiesService bagsiesS.Service,
) {
	authMiddleware := middlewares.NewAuth(authService, usersService)

	authController := auth.New(authService, usersService)
	formsController := form.NewController(formsService)
	bagsiesController := bagsies.New(bagsiesService, authService, usersService, authMiddleware)

	server.WithMiddlewares(initCORSMiddleware())
	server.WithControllers(authController, formsController, bagsiesController)

	initSwagger(server, swaggerHost, swaggerScheme)
}

func initSwagger(server *coreHttp.Server, swaggerHost, swaggerScheme string) {
	if version.GetVersion() != enum.EnvironmentProd {
		docs.SwaggerInfo.Host = swaggerHost
		docs.SwaggerInfo.Schemes = []string{swaggerScheme}
		server.WithControllers(swagger.New(swaggerScheme, swaggerHost))
		return
	}
	log.Warn(context.Background(), "version is not supported", log.String("version", version.GetVersion().String()))
}

func initCORSMiddleware() coreHttp.Middleware {
	corsMiddleware := coreHttp.NewCORSMiddleware().
		WithMethods(
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		).
		WithCredentials(true).
		WithOrigins(
			"*",
		).
		WithHeaders(
			"Accept",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
		)
	return corsMiddleware
}
