package analytics

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	coreHTTP "github.com/Rasikrr/core/http"
)

// orgContext извлекает OrgContext или отправляет 401.
func orgContext(w http.ResponseWriter, r *http.Request) (*access.OrgContext, bool) {
	ctx := r.Context()
	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return nil, false
	}
	return orgCtx, true
}
