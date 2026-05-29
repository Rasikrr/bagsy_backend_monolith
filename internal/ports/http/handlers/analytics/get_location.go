package analytics

import (
	"net/http"

	domainAnalytics "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/analytics"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// getLocation handles GET /api/v1/analytics/locations/{locationID}.
//
// @Summary      Drill-down по локации
// @Description  Сводка аналитики в скоупе одной локации. Только для Network-owner.
// @Tags         analytics
// @Produce      json
// @Param        locationID  path      string  true  "ID локации (UUID)"
// @Param        from        query     string  true  "Начало периода (YYYY-MM-DD)"
// @Param        to          query     string  true  "Конец периода (YYYY-MM-DD)"
// @Success      200  {object}  overviewResponse
// @Failure      401  {object}  httputil.errorResponse
// @Failure      403  {object}  httputil.errorResponse
// @Failure      404  {object}  httputil.errorResponse
// @Failure      422  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/analytics/locations/{locationID} [get]
func (h *Handler) getLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := orgContext(w, r)
	if !ok {
		return
	}

	locationID, err := uuid.Parse(chi.URLParam(r, "locationID"))
	if err != nil {
		httputil.SendError(ctx, w, domainAnalytics.ErrNotFound, analyticsErrors)
		return
	}

	params, err := parsePeriodParams(r)
	if err != nil {
		httputil.SendError(ctx, w, err, analyticsErrors)
		return
	}

	report, err := h.analyticsUC.GetLocation(ctx, orgCtx, locationID, params.From, params.To)
	if err != nil {
		httputil.SendError(ctx, w, err, analyticsErrors)
		return
	}

	setCacheHeaders(w)
	coreHTTP.SendData(ctx, w, toOverviewResponse(report), http.StatusOK)
}
