// nolint
package http

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/calendar"
	masterservicesH "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/master_services"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/media"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/networks"
	pointcategories "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/point_categories"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/points"
	servicecategories "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/service_categories"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/services"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/users"
	masterServicesS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/master_services"
	mediaS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/media"
	networksS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/networks"
	pointCategoriesS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/point_categories"
	pointsS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/points"
	serviceCategoriesS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/service_categories"
	servicesS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/services"
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
	servicesService *servicesS.Service,
	mediaService *mediaS.Service,
	pointCategoriesService *pointCategoriesS.Service,
	serviceCategoriesService *serviceCategoriesS.Service,
	masterServicesService *masterServicesS.Service,
) {
	authMiddleware := middlewares.NewAuth(authService)
	rateLimiterFactory := middlewares.NewRateLimiterFactory(redis)

	authController := auth.New(authService, authMiddleware, rateLimiterFactory)
	formsController := form.New(formsService)
	usersController := users.New(usersService, authMiddleware)
	bagsiesController := bagsies.New(bagsiesService, authMiddleware)
	pointsController := points.New(pointsService, authMiddleware)
	networksController := networks.New(networksService, pointsService, authMiddleware)
	servicesController := services.New(servicesService, authMiddleware)
	mediaController := media.New(mediaService, authMiddleware)
	calendarController := calendar.New(bagsiesService, authMiddleware)
	pointCategoriesController := pointcategories.New(pointCategoriesService)
	serviceCategoriesController := servicecategories.New(serviceCategoriesService)
	masterServicesController := masterservicesH.New(masterServicesService, authMiddleware)

	server.WithMiddlewares(initCORSMiddleware())
	server.WithControllers(
		authController,
		formsController,
		usersController,
		bagsiesController,
		pointsController,
		networksController,
		servicesController,
		mediaController,
		calendarController,
		pointCategoriesController,
		serviceCategoriesController,
		masterServicesController,
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
