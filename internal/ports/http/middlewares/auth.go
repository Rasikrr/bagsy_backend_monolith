package middlewares

import (
	"context"
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/actor"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
)

// Определен здесь согласно DIP - интерфейс определяется где используется
type authService interface {
	VerifyAccessToken(ctx context.Context, tokenStr string) (*actor.Actor, error)
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

		act, err := a.authService.VerifyAccessToken(ctx, token)
		if err != nil {
			errors.HandleError(ctx, w, err)
			return
		}

		ctx = actor.SetActor(ctx, act)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Auth) RequireRole(roles ...user.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return a.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			act, err := actor.GetActor(ctx)
			if err != nil {
				errors.HandleError(ctx, w, err)
				return
			}

			if !act.Role().OneOf(roles...) {
				errors.HandleError(ctx, w, domainErr.NewForbiddenError("insufficient permissions"))
				return
			}

			next.ServeHTTP(w, r)
		}))
	}
}

func (a *Auth) AuthorizeManagement() func(http.Handler) http.Handler {
	return a.RequireRole(
		user.RoleManager,
		user.RoleNetManager,
		user.RoleSelfOwner,
	)
}

func (a *Auth) AuthorizeNetManagement() func(http.Handler) http.Handler {
	return a.RequireRole(
		user.RoleNetManager,
		user.RoleSelfOwner,
	)
}

func (a *Auth) AuthorizeWorkers() func(http.Handler) http.Handler {
	return a.RequireRole(
		user.RoleStaff,
		user.RoleManager,
		user.RoleNetManager,
		user.RoleSelfOwner,
	)
}
