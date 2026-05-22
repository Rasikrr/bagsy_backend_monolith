package billing

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
)

// undoCancellation handles POST /api/v1/subscription/undo-cancel.
//
// @Summary      Отмена запроса на отмену подписки
// @Description  Отменяет ранее запрошенную отмену подписки. Подписка продолжит автопродление. Только Owner.
// @Tags         billing
// @Success      204
// @Failure      401  {object}  httputil.errorResponse
// @Failure      403  {object}  httputil.errorResponse
// @Failure      422  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/subscription/undo-cancel [post]
func (h *Handler) undoCancellation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	if err := h.billingUC.UndoCancellation(ctx, orgCtx); err != nil {
		httputil.SendError(ctx, w, err, billingErrors)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
