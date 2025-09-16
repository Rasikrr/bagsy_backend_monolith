package http

import (
	"github.com/Rasikrr/bugsy_backend_monolith/internal/ports/http/handlers/auth"
	authS "github.com/Rasikrr/bugsy_backend_monolith/internal/services/auth"
	"github.com/Rasikrr/core/http"
)

func NewServer(
	server *http.Server,
	authService authS.Service,
) {
	authController := auth.New(authService)
	server.WithControllers(authController)
}
