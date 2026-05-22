package billing

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/billing"
	coreHTTP "github.com/Rasikrr/core/http"
)

// activate handles POST /api/v1/subscription/activate.
//
// @Summary      Активация подписки
// @Description  Активирует подписку с выбранным циклом оплаты (monthly/annual). Только Owner. Допустимо из статусов: trial, past_due, suspended.
// @Tags         billing
// @Accept       json
// @Param        body  body  activateRequest  true  "Цикл оплаты"
// @Success      204
// @Failure      400   {object}  httputil.errorResponse
// @Failure      401   {object}  httputil.errorResponse
// @Failure      403   {object}  httputil.errorResponse
// @Failure      422   {object}  httputil.errorResponse
// @Failure      500   {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/subscription/activate [post]
func (h *Handler) activate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	var req activateRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	input := uc.ActivateInput{
		Cycle: req.Cycle,
	}

	if err := h.billingUC.Activate(ctx, orgCtx, input); err != nil {
		httputil.SendError(ctx, w, err, billingErrors)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
