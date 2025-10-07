package middlewares

import (
	"errors"
	"net/http"

	"github.com/Rasikrr/core/log"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/session"
	"github.com/Rasikrr/core/api"
)

type AuthMiddleware struct {
	authService  auth.Service
	usersService users.Service
}

func NewAuth(
	authService auth.Service,
	usersService users.Service,
) AuthMiddleware {
	return AuthMiddleware{
		authService:  authService,
		usersService: usersService,
	}
}

// nolint: nonamedreturns
func (a *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := session.GetAuthHeader(r)
		log.Debugf(r.Context(), "token in middleware = %v", token)
		if err != nil {
			api.SendError(w, err)
			return
		}

		ctx := r.Context()
		payload, err := a.authService.GetAuthTokenPayload(ctx, token)
		if err != nil {
			api.SendError(w, err)
			return
		}
		if payload.IsRefresh() {
			api.SendError(w, errors.New("refresh token is not allowed"))
			return
		}
		ses, err := payload.ToSession()
		if err != nil {
			api.SendError(w, err)
			return
		}
		log.Infof(ctx, "set session %+v", ses)
		ctx = session.SetSession(ctx, ses)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
