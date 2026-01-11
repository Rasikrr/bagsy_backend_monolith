// nolint
package http

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/media"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/networks"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/points"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/users"
	mediaS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/media"
	networksS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/networks"
	pointsS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/points"
	usersS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
	"github.com/Rasikrr/core/cache/redis"
	"github.com/Rasikrr/core/environment"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/go-chi/cors"

	docs "github.com/Rasikrr/bagsy_backend_monolith/docs/swagger"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/bagsies"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/form"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/swagger"
	authS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"
	bagsiesS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/bagsies"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/forms"
	"github.com/Rasikrr/core/enum"
	coreHttp "github.com/Rasikrr/core/http"
)

// @title Bagsy API
// @version 1.0
// @description API для приложения Bagsy
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@bagsies.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewServer(
	server *coreHttp.Server,
	redis *redis.Client,
	swaggerHost, swaggerScheme string,
	authService *authS.Service,
	formsService *forms.Service,
	usersService *usersS.Service,
	bagsiesService *bagsiesS.Service,
	pointsService *pointsS.Service,
	networksService *networksS.Service,
	mediaService *mediaS.Service,
) {
	authMiddleware := middlewares.NewAuth(authService)
	rateLimiterFactory := middlewares.NewRateLimiterFactory(redis)

	authController := auth.New(authService, authMiddleware, rateLimiterFactory)
	formsController := form.New(formsService)
	usersController := users.New(usersService, authMiddleware)
	bagsiesController := bagsies.New(bagsiesService, authMiddleware)
	pointsController := points.New(pointsService, authMiddleware)
	networksController := networks.New(networksService, pointsService, authMiddleware)
	mediaController := media.New(mediaService, authMiddleware)

	server.WithMiddlewares(initCORSMiddleware())
	server.WithControllers(
		authController,
		formsController,
		usersController,
		bagsiesController,
		pointsController,
		networksController,
		mediaController,
	)

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
