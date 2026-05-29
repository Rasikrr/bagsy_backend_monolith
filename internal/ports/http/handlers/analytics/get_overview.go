package analytics

import (
	"net/http"

	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
)

// getOverview handles GET /api/v1/analytics/overview.
//
// @Summary      Сводка аналитики
// @Description  KPI, выручка по дням, топы мастеров/услуг, воронка, heatmap и инсайты за период.
// @Tags         analytics
// @Produce      json
// @Param        from        query     string  true   "Начало периода (YYYY-MM-DD)"
// @Param        to          query     string  true   "Конец периода (YYYY-MM-DD)"
// @Param        location_id query     string  false  "ID локации (UUID)"
// @Success      200  {object}  overviewResponse
// @Failure      401  {object}  httputil.errorResponse
// @Failure      403  {object}  httputil.errorResponse
// @Failure      422  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/analytics/overview [get]
func (h *Handler) getOverview(w http.ResponseWriter, r *http.Request) {
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

	report, err := h.analyticsUC.GetOverview(ctx, orgCtx, params.From, params.To, params.LocationID)
	if err != nil {
		httputil.SendError(ctx, w, err, analyticsErrors)
		return
	}

	setCacheHeaders(w)
	coreHTTP.SendData(ctx, w, toOverviewResponse(report), http.StatusOK)
}
