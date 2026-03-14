package schedule

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// deleteEmployeeSchedule handles DELETE /api/v1/employee-schedules/{employeeID}.
//
// @Summary      Удаление расписания сотрудника
// @Description  Удаляет все слоты расписания сотрудника за указанный период. Owner — любого, Manager — staff в своей локации.
// @Tags         schedule
// @Param        employeeID  path      string  true  "ID сотрудника"
// @Param        start       query     string  true  "Начало периода (YYYY-MM-DD)"
// @Param        end         query     string  true  "Конец периода (YYYY-MM-DD)"
// @Success      204  "No Content"
// @Failure      400  {object}  httputil.errorResponse
// @Failure      401  {object}  httputil.errorResponse
// @Failure      403  {object}  httputil.errorResponse
// @Failure      404  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/employee-schedules/{employeeID} [delete]
func (h *Handler) deleteEmployeeSchedule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	employeeID, err := uuid.Parse(chi.URLParam(r, "employeeID"))
	if err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	start, end, err := parseDateRange(r)
	if err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	if err = h.scheduleUC.DeleteEmployeeSchedule(ctx, orgCtx, employeeID, start, end); err != nil {
		httputil.SendError(ctx, w, err, scheduleErrors)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
