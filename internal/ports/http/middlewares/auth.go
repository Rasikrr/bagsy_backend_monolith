package middlewares

import (
	"github.com/Rasikrr/bugsy_backend_monolith/internal/services/auth"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/services/users"
)

type AuthMiddleware struct {
	authService  auth.Service
	usersService users.Service
}

func NewAuth(
	authService auth.Service,
	usersService users.Service,
) *AuthMiddleware {
	return &AuthMiddleware{
		authService:  authService,
		usersService: usersService,
	}
}
