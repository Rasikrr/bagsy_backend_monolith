package middlewares

import (
	"context"
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
)

type authUseCase interface {
	VerifyAccessToken(ctx context.Context, tokenStr string) (*auth.Token, error)
}

type Auth struct {
	authService authUseCase
}

func NewAuth(authService authUseCase) *Auth {
	return &Auth{
		authService: authService,
	}
}

func (a *Auth) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token, err := util.GetAuthHeader(r)
		if err != nil {
			util.SendBadRequest(ctx, w, err)
			return
		}

		tokenInfo, err := a.authService.VerifyAccessToken(ctx, token)
		if err != nil {
			util.SendBadRequest(ctx, w, err)
			return
		}

		ctx = access.WithToken(ctx, tokenInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
