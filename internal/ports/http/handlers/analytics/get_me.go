package analytics

import (
	"net/http"

	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
)

// getMe handles GET /api/v1/analytics/me.
//
// @Summary      Личная аналитика сотрудника
// @Description  KPI, выручка по дням, топ услуг, heatmap и разбивка клиентов текущего сотрудника.
// @Tags         analytics
// @Produce      json
// @Param        from  query     string  true  "Начало периода (YYYY-MM-DD)"
// @Param        to    query     string  true  "Конец периода (YYYY-MM-DD)"
// @Success      200  {object}  meResponse
// @Failure      401  {object}  httputil.errorResponse
// @Failure      422  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/analytics/me [get]
func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := orgContext(w, r)
	if !ok {
		return
	}

	params, err := parsePeriodParams(r)
	if err != nil {
		httputil.SendError(ctx, w, err, analyticsErrors)
		return
	}

	report, err := h.analyticsUC.GetMe(ctx, orgCtx, params.From, params.To)
	if err != nil {
		httputil.SendError(ctx, w, err, analyticsErrors)
		return
	}

	setCacheHeaders(w)
	coreHTTP.SendData(ctx, w, toMeResponse(report), http.StatusOK)
}
