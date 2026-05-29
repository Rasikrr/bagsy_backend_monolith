package analytics

import (
	"net/http"

	domainAnalytics "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/analytics"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// getStaffDetail handles GET /api/v1/analytics/staff/{employeeID}.
//
// @Summary      Drill-down по мастеру
// @Description  KPI, выручка по дням, топ услуг, почасовая нагрузка и разбивка клиентов мастера.
// @Tags         analytics
// @Produce      json
// @Param        employeeID  path      string  true   "ID мастера (UUID)"
// @Param        from        query     string  true   "Начало периода (YYYY-MM-DD)"
// @Param        to          query     string  true   "Конец периода (YYYY-MM-DD)"
// @Success      200  {object}  staffDetailResponse
// @Failure      401  {object}  httputil.errorResponse
// @Failure      403  {object}  httputil.errorResponse
// @Failure      404  {object}  httputil.errorResponse
// @Failure      422  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/analytics/staff/{employeeID} [get]
func (h *Handler) getStaffDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := orgContext(w, r)
	if !ok {
		return
	}

	employeeID, err := uuid.Parse(chi.URLParam(r, "employeeID"))
	if err != nil {
		httputil.SendError(ctx, w, domainAnalytics.ErrNotFound, analyticsErrors)
		return
	}

	params, err := parsePeriodParams(r)
	if err != nil {
		httputil.SendError(ctx, w, err, analyticsErrors)
		return
	}

	report, err := h.analyticsUC.GetStaffDetail(ctx, orgCtx, employeeID, params.From, params.To)
	if err != nil {
		httputil.SendError(ctx, w, err, analyticsErrors)
		return
	}

	setCacheHeaders(w)
	coreHTTP.SendData(ctx, w, toStaffDetailResponse(report), http.StatusOK)
}
