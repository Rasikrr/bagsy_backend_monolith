package schedule

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	domainSchedule "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var errStartEndRequired = errors.New("start and end query params are required")

// getLocationSchedule handles GET /api/v1/location-schedules/{locationID}.
//
// @Summary      Получение расписания локации
// @Description  Возвращает слоты расписания локации за указанный период.
// @Tags         schedule
// @Produce      json
// @Param        locationID  path      string  true  "ID локации"
// @Param        start       query     string  true  "Начало периода (YYYY-MM-DD)"
// @Param        end         query     string  true  "Конец периода (YYYY-MM-DD)"
// @Success      200  {object}  getScheduleResponse
// @Failure      400  {object}  httputil.errorResponse
// @Failure      401  {object}  httputil.errorResponse
// @Failure      403  {object}  httputil.errorResponse
// @Failure      404  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/location-schedules/{locationID} [get]
func (h *Handler) getLocationSchedule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	locationID, err := uuid.Parse(chi.URLParam(r, "locationID"))
	if err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	start, end, err := parseDateRange(r)
	if err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	slots, err := h.scheduleUC.GetLocationSchedule(ctx, orgCtx, locationID, start, end)
	if err != nil {
		httputil.SendError(ctx, w, err, scheduleErrors)
		return
	}

	sendSlotResponse(ctx, w, toLocationSlotResponses(slots))
}

func toLocationSlotResponses(slots []*domainSchedule.LocationScheduleSlot) []slotResponse {
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

func sendSlotResponse(ctx context.Context, w http.ResponseWriter, slots []slotResponse) {
	coreHTTP.SendData(ctx, w, getScheduleResponse{Slots: slots}, http.StatusOK)
}

func parseDateRange(r *http.Request) (time.Time, time.Time, error) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	if startStr == "" || endStr == "" {
		return time.Time{}, time.Time{}, errStartEndRequired
	}

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return start, end, nil
}
