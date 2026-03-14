package schedule

import (
	"net/http"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/schedule"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// setEmployeeSchedule handles PUT /api/v1/employee-schedules/{employeeID}.
//
// @Summary      Установка расписания сотрудника
// @Description  Заменяет все слоты сотрудника за указанный период новыми. Owner — любого, Manager — staff в своей локации.
// @Tags         schedule
// @Accept       json
// @Param        employeeID  path      string              true  "ID сотрудника"
// @Param        body        body      setScheduleRequest  true  "Период и слоты расписания"
// @Success      204  "No Content"
// @Failure      400  {object}  httputil.errorResponse
// @Failure      401  {object}  httputil.errorResponse
// @Failure      403  {object}  httputil.errorResponse
// @Failure      404  {object}  httputil.errorResponse
// @Failure      422  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/employee-schedules/{employeeID} [put]
func (h *Handler) setEmployeeSchedule(w http.ResponseWriter, r *http.Request) {
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

	var req setScheduleRequest
	if err = coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	start, err := time.Parse("2006-01-02", req.Start)
	if err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	end, err := time.Parse("2006-01-02", req.End)
	if err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	slots, err := toSlotInputs(req.Slots)
	if err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	input := uc.SetEmployeeScheduleInput{
		EmployeeID: employeeID,
		Start:      start,
		End:        end,
		Slots:      slots,
	}

	if err = h.scheduleUC.SetEmployeeSchedule(ctx, orgCtx, input); err != nil {
		httputil.SendError(ctx, w, err, scheduleErrors)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
