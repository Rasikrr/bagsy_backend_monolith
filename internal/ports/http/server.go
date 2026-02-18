package http

import (
	docs "github.com/Rasikrr/bagsy_backend_monolith/docs/swagger"
	authC "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/swagger"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/auth"
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
) {
	server.WithMiddlewares(initCORSMiddleware())
	initSwagger(server, swaggerHost, swaggerScheme)

	authHandler := authC.New(registerOwnerUseCase, authUseCase, resetPasswordUseCase)

	server.WithControllers(
		authHandler,
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
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedOrigins: []string{"*"},
		},
	)
	return corsMiddleware
}
