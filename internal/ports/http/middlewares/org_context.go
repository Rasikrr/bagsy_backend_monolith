package middlewares

import (
	"context"
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/google/uuid"
)

type orgContextRepository interface {
	GetOrgContext(ctx context.Context, employeeID uuid.UUID) (*access.OrgContext, error)
}

type OrgContext struct {
	repo orgContextRepository
}

func NewOrgContext(repo orgContextRepository) *OrgContext {
	return &OrgContext{
		repo: repo,
	}
}

func (m *OrgContext) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token, ok := access.TokenFromContext(ctx)
		if !ok {
			coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
			return
		}

		orgCtx, err := m.repo.GetOrgContext(ctx, token.UserID)
		if err != nil {
			coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
			return
		}

		if !orgCtx.Organization.Active {
			coreHTTP.SendData(ctx, w, map[string]string{"error": "organization_inactive"}, http.StatusForbidden)
			return
		}

		ctx = access.WithOrgContext(ctx, orgCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
