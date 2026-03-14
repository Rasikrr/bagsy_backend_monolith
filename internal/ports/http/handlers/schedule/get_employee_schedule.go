package schedule

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	domainSchedule "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// getEmployeeSchedule handles GET /api/v1/employee-schedules/{employeeID}.
//
// @Summary      Получение расписания сотрудника
// @Description  Возвращает слоты расписания сотрудника за указанный период. Owner — любого, Manager — своей локации, Staff — только своё.
// @Tags         schedule
// @Produce      json
// @Param        employeeID  path      string  true  "ID сотрудника"
// @Param        start       query     string  true  "Начало периода (YYYY-MM-DD)"
// @Param        end         query     string  true  "Конец периода (YYYY-MM-DD)"
// @Success      200  {object}  getScheduleResponse
// @Failure      400  {object}  httputil.errorResponse
// @Failure      401  {object}  httputil.errorResponse
// @Failure      403  {object}  httputil.errorResponse
// @Failure      404  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/employee-schedules/{employeeID} [get]
func (h *Handler) getEmployeeSchedule(w http.ResponseWriter, r *http.Request) {
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

	slots, err := h.scheduleUC.GetEmployeeSchedule(ctx, orgCtx, employeeID, start, end)
	if err != nil {
		httputil.SendError(ctx, w, err, scheduleErrors)
		return
	}

	sendSlotResponse(ctx, w, toEmployeeSlotResponses(slots))
}

func toEmployeeSlotResponses(slots []*domainSchedule.EmployeeScheduleSlot) []slotResponse {
	resp := make([]slotResponse, 0, len(slots))
	for _, s := range slots {
		resp = append(resp, slotResponse{
			ID:        s.ID.String(),
			Date:      s.Date.Format("2006-01-02"),
			Type:      s.Type.String(),
			StartTime: s.StartTime.Format("15:04"),
			EndTime:   s.EndTime.Format("15:04"),
		})
	}
	return resp
}
