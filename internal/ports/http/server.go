// nolint
package http

import (
	"github.com/Rasikrr/core/environment"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/go-chi/cors"

	docs "github.com/Rasikrr/bagsy_backend_monolith/docs/swagger"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/form"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/swagger"
	authS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/forms"
	"github.com/Rasikrr/core/enum"
	coreHttp "github.com/Rasikrr/core/http"
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
	authService *authS.Service,
	formsService *forms.Service,
) {
	authMiddleware := middlewares.NewAuth(authService)

	authController := auth.New(authService, authMiddleware)
	formsController := form.New(formsService)

	server.WithMiddlewares(initCORSMiddleware())
	server.WithControllers(authController, formsController)

	initSwagger(server, swaggerHost, swaggerScheme)
}

func initSwagger(server *coreHttp.Server, swaggerHost, swaggerScheme string) {
	if environment.GetEnv() != enum.EnvironmentProd {
		docs.SwaggerInfo.Host = swaggerHost
		docs.SwaggerInfo.Schemes = []string{swaggerScheme}
		server.WithControllers(swagger.New(swaggerScheme, swaggerHost))
		return
	}
}

func initCORSMiddleware() coreHttp.Middleware {
	corsMiddleware := coreHttp.NewCORSMiddleware(
		cors.Options{
			AllowCredentials: false,
			AllowedHeaders: []string{
				"Accept",
				"Content-Type",
				"Content-Length",
				"Accept-Encoding",
				"X-CSRF-Token",
				"Authorization",
				coreHttp.TraceIDHeader,
			},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedOrigins: []string{"*"},
		},
	)
	return corsMiddleware
}
