package middlewares

import (
	"context"
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/session"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
)

// Определен здесь согласно DIP - интерфейс определяется где используется
type authService interface {
	VerifyAccessToken(ctx context.Context, tokenStr string) (*session.Session, error)
}

type AuthMiddleware struct {
	authService authService
}

func NewAuth(authService authService) AuthMiddleware {
	return AuthMiddleware{
		authService: authService,
	}
}

// Handle базовый middleware для проверки токена
func (a *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token, err := httputil.GetAuthHeader(r)
		if err != nil {
			errors.HandleError(ctx, w, err)
			return
		}

		ses, err := a.authService.VerifyAccessToken(ctx, token)
		if err != nil {
			errors.HandleError(ctx, w, err)
			return
		}
		// Устанавливаем сессию в контекст
		ctx = session.SetSession(ctx, ses)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// RequireRole проверяет что пользователь имеет одну из указанных ролей
func (a *AuthMiddleware) RequireRole(roles ...enum.Role) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return a.Handle(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ses, err := session.GetSession(ctx)
			if err != nil {
				errors.HandleError(ctx, w, err)
				return
			}

			if !ses.Role().OneOf(roles...) {
				errors.HandleError(ctx, w, domainErr.NewForbiddenError("insufficient permissions"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
