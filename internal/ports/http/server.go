// nolint
package http

import (
	"github.com/Rasikrr/bugsy_backend_monolith/internal/ports/http/handlers/auth"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/ports/http/handlers/swagger"
	authS "github.com/Rasikrr/bugsy_backend_monolith/internal/services/auth"
	usersS "github.com/Rasikrr/bugsy_backend_monolith/internal/services/users"
	"github.com/Rasikrr/core/enum"
	"github.com/Rasikrr/core/http"
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
	authService authS.Service,
	usersService usersS.Service,
) {
	authController := auth.New(authService, usersService)
	server.WithControllers(authController)

	initSwagger(server)
}

func initSwagger(server *http.Server) {
	if version.GetVersion() == enum.EnvironmentDev {
		server.WithControllers(swagger.New())
	}
}
