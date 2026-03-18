package http

import (
	docs "github.com/Rasikrr/bagsy_backend_monolith/docs/swagger"
	authC "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/auth"
	bookingC "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/booking"
	catalogC "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/catalog"
	employeeC "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/employee"
	locationC "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/location"
	mediaC "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/media"
	scheduleC "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/schedule"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/swagger"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/auth"
	bookingUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/booking"
	catalogUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/catalog"
	employeeUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/employee"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/invite"
	locationUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/location"
	mediaUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/media"
	scheduleUC "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/schedule"
	"github.com/Rasikrr/core/enum"
	"github.com/Rasikrr/core/environment"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/cors"
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
	server *coreHTTP.Server,
	swaggerHost, swaggerScheme string,
	registerOwnerUseCase *auth.RegisterOwnerUseCase,
	authUseCase *auth.UseCase,
	resetPasswordUseCase *auth.ResetPasswordUseCase,
	inviteUseCase *invite.UseCase,
	employeeUseCase *employeeUC.UseCase,
	accessRepo *access.Repository,
	createLocationUC *locationUC.UseCase,
	catalogUseCase *catalogUC.UseCase,
	bookingUseCase *bookingUC.UseCase,
	mediaUseCase *mediaUC.UseCase,
	scheduleUseCase *scheduleUC.UseCase,
) {
	server.WithMiddlewares(initCORSMiddleware())
	initSwagger(server, swaggerHost, swaggerScheme)

	authMiddleware := middlewares.NewAuth(authUseCase)
	orgContextMiddleware := middlewares.NewOrgContext(accessRepo)

	authHandler := authC.New(registerOwnerUseCase, authUseCase, resetPasswordUseCase)
	locationHandler := locationC.New(createLocationUC, authMiddleware, orgContextMiddleware)
	employeeHandler := employeeC.New(inviteUseCase, employeeUseCase, catalogUseCase, authMiddleware, orgContextMiddleware)
	catalogHandler := catalogC.New(catalogUseCase, authMiddleware, orgContextMiddleware)
	bookingHandler := bookingC.New(bookingUseCase, authMiddleware, orgContextMiddleware)
	mediaHandler := mediaC.New(mediaUseCase, authMiddleware)
	scheduleHandler := scheduleC.New(scheduleUseCase, authMiddleware, orgContextMiddleware)

	server.WithControllers(
		authHandler,
		locationHandler,
		employeeHandler,
		catalogHandler,
		bookingHandler,
		mediaHandler,
		scheduleHandler,
	)
}

func initSwagger(server *coreHTTP.Server, swaggerHost, swaggerScheme string) {
	if environment.GetEnv() != enum.EnvironmentProd {
		docs.SwaggerInfo.Host = swaggerHost
		docs.SwaggerInfo.Schemes = []string{swaggerScheme}
		server.WithControllers(
			swagger.New(swaggerScheme, swaggerHost),
		)
		return
	}
}

func initCORSMiddleware() coreHTTP.Middleware {
	corsMiddleware := coreHTTP.NewCORSMiddleware(
		cors.Options{
			AllowCredentials: false,
			AllowedHeaders: []string{
				"Accept",
				"Content-Type",
				"Content-Length",
				"Accept-Encoding",
				"X-CSRF-Token",
				"Authorization",
				coreHTTP.TraceIDHeader,
			},
			AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedOrigins: []string{"*"},
		},
	)
	return corsMiddleware
}
