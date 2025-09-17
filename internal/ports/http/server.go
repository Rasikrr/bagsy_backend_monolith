package http

import (
	"github.com/Rasikrr/bugsy_backend_monolith/internal/ports/http/handlers/auth"
	authS "github.com/Rasikrr/bugsy_backend_monolith/internal/services/auth"
	usersS "github.com/Rasikrr/bugsy_backend_monolith/internal/services/users"
	"github.com/Rasikrr/core/http"
)

func NewServer(
	server *http.Server,
	authService authS.Service,
	usersService usersS.Service,
) {
	authController := auth.New(authService, usersService)
	server.WithControllers(authController)
}
