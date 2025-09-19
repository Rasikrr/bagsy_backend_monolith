// nolint
package http

import (
	"context"

	docs "github.com/Rasikrr/bugsy_backend_monolith/docs/swagger"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/ports/http/handlers/auth"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/ports/http/handlers/swagger"
	authS "github.com/Rasikrr/bugsy_backend_monolith/internal/services/auth"
	usersS "github.com/Rasikrr/bugsy_backend_monolith/internal/services/users"
	"github.com/Rasikrr/core/enum"
	"github.com/Rasikrr/core/http"
	"github.com/Rasikrr/core/log"
	"github.com/Rasikrr/core/version"
)

// @title Bugsy API
// @version 1.0
// @description API для системы Bugsy
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewServer(
	server *http.Server,
	swaggerHost, swaggerScheme string,
	authService authS.Service,
	usersService usersS.Service,
) {
	authController := auth.New(authService, usersService)
	server.WithControllers(authController)

	initSwagger(server, swaggerHost, swaggerScheme)
}

func initSwagger(server *http.Server, swaggerHost, swaggerScheme string) {
	if version.GetVersion() != enum.EnvironmentProd {
		docs.SwaggerInfo.Host = swaggerHost
		docs.SwaggerInfo.Schemes = []string{swaggerScheme}
		server.WithControllers(swagger.New())
		return
	}
	log.Warn(context.Background(), "version is not supported", log.String("version", version.GetVersion().String()))
}
