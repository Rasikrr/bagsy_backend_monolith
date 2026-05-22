package billing

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
)

// requestCancellation handles POST /api/v1/subscription/cancel.
//
// @Summary      Запрос отмены подписки
// @Description  Помечает подписку для отмены в конце текущего оплаченного периода. Только Owner. Допустимо из статуса: active.
// @Tags         billing
// @Success      204
// @Failure      401  {object}  httputil.errorResponse
// @Failure      403  {object}  httputil.errorResponse
// @Failure      409  {object}  httputil.errorResponse
// @Failure      422  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/subscription/cancel [post]
func (h *Handler) requestCancellation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	if err := h.billingUC.RequestCancellation(ctx, orgCtx); err != nil {
		httputil.SendError(ctx, w, err, billingErrors)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
