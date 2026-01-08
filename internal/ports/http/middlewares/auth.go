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

type Auth struct {
	authService authService
}

func NewAuth(authService authService) *Auth {
	return &Auth{
		authService: authService,
	}
}

func (a *Auth) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx = session.SetSession(ctx, ses)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Auth) RequireRole(roles ...enum.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return a.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		}))
	}
}

func (a *Auth) AuthorizeManagement() func(http.Handler) http.Handler {
	return a.RequireRole(
		enum.RoleManager,
		enum.RoleNetManager,
		enum.RoleSelfOwner,
	)
}

func (a *Auth) AuthorizeNetManagement() func(http.Handler) http.Handler {
	return a.RequireRole(
		enum.RoleNetManager,
		enum.RoleSelfOwner,
	)
}

func (a *Auth) AuthorizeWorkers() func(http.Handler) http.Handler {
	return a.RequireRole(
		enum.RoleStaff,
		enum.RoleManager,
		enum.RoleNetManager,
		enum.RoleSelfOwner,
	)
}
