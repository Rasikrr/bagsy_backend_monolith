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

// setLocationSchedule handles PUT /api/v1/location-schedules/{locationID}.
//
// @Summary      Установка расписания локации
// @Description  Заменяет все слоты локации за указанный период новыми. Удаляет старые, вставляет новые в одной транзакции.
// @Tags         schedule
// @Accept       json
// @Param        locationID  path      string              true  "ID локации"
// @Param        body        body      setScheduleRequest  true  "Период и слоты расписания"
// @Success      204  "No Content"
// @Failure      400  {object}  httputil.errorResponse
// @Failure      401  {object}  httputil.errorResponse
// @Failure      403  {object}  httputil.errorResponse
// @Failure      404  {object}  httputil.errorResponse
// @Failure      422  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/location-schedules/{locationID} [put]
func (h *Handler) setLocationSchedule(w http.ResponseWriter, r *http.Request) {
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

	input := uc.SetLocationScheduleInput{
		LocationID: locationID,
		Start:      start,
		End:        end,
		Slots:      slots,
	}

	if err = h.scheduleUC.SetLocationSchedule(ctx, orgCtx, input); err != nil {
		httputil.SendError(ctx, w, err, scheduleErrors)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toSlotInputs(reqs []slotRequest) ([]uc.SlotInput, error) {
	slots := make([]uc.SlotInput, 0, len(reqs))
	for _, r := range reqs {
		date, err := time.Parse("2006-01-02", r.Date)
		if err != nil {
			return nil, err
		}

		startTime, err := time.Parse("15:04", r.StartTime)
		if err != nil {
			return nil, err
		}

		endTime, err := time.Parse("15:04", r.EndTime)
		if err != nil {
			return nil, err
		}

		slots = append(slots, uc.SlotInput{
			Date:      date,
			Type:      r.Type,
			StartTime: startTime,
			EndTime:   endTime,
		})
	}
	return slots, nil
}
