package http

import (
	"github.com/Rasikrr/bugsy_backend_monolith/internal/ports/http/handlers/auth"
	"github.com/Rasikrr/core/http"
)

func NewServer(
	server *http.Server,
) {
	authController := auth.New()
	server.WithControllers(authController)
}
