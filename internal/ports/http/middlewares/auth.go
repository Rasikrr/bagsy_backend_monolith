package middlewares

import (
	"errors"
	"net/http"

	"github.com/Rasikrr/core/log"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/session"
	"github.com/Rasikrr/core/api"
	coreErr "github.com/Rasikrr/core/errors"
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

// RequireRole проверяет что пользователь имеет одну из указанных ролей
func (a *AuthMiddleware) RequireRole(roles ...enum.Role) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return a.Handle(func(w http.ResponseWriter, r *http.Request) {
			ses, err := session.GetSession(r.Context())
			if err != nil {
				api.SendError(w, err)
				return
			}

			if !ses.Role.OneOf(roles...) {
				api.SendError(w, coreErr.ErrForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireMinRole проверяет что пользователь имеет минимальный уровень роли
func (a *AuthMiddleware) RequireMinRole(minRole enum.Role) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return a.Handle(func(w http.ResponseWriter, r *http.Request) {
			ses, err := session.GetSession(r.Context())
			if err != nil {
				api.SendError(w, err)
				return
			}

			if ses.Role < minRole {
				api.SendError(w, coreErr.ErrForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
