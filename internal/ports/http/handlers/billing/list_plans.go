package billing

import (
	"net/http"

	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
)

// listPlans handles GET /api/v1/plans.
//
// @Summary      Список тарифных планов
// @Description  Возвращает все активные тарифные планы с ценами и лимитами. Публичный эндпоинт.
// @Tags         billing
// @Produce      json
// @Success      200  {object}  listPlansResponse
// @Failure      500  {object}  httputil.errorResponse
// @Router       /api/v1/plans [get]
func (h *Handler) listPlans(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	plans, err := h.billingUC.ListPlans(ctx)
	if err != nil {
		httputil.SendError(ctx, w, err, billingErrors)
		return
	}

	resp := make([]planResponse, 0, len(plans))
	for _, p := range plans {
		resp = append(resp, toPlanResponse(p))
	}

	coreHTTP.SendData(ctx, w, listPlansResponse{Plans: resp}, http.StatusOK)
}
