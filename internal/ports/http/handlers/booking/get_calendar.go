package booking

import (
	"net/http"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/booking"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/google/uuid"
)

func (h *Handler) getCalendar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	q := r.URL.Query()

	from, err := time.Parse("2006-01-02", q.Get("from"))
	if err != nil {
		util.SendBadRequest(ctx, w, err)
		return
	}

	to, err := time.Parse("2006-01-02", q.Get("to"))
	if err != nil {
		util.SendBadRequest(ctx, w, err)
		return
	}

	input := uc.GetCalendarInput{
		StartDate:        from,
		EndDate:          to,
		IncludeCancelled: q.Get("include_cancelled") == "true",
	}

	if locID := q.Get("location_id"); locID != "" {
		parsed, err := uuid.Parse(locID)
		if err != nil {
			util.SendBadRequest(ctx, w, err)
			return
		}
		input.LocationID = &parsed
	}

	if empID := q.Get("employee_id"); empID != "" {
		parsed, err := uuid.Parse(empID)
		if err != nil {
			util.SendBadRequest(ctx, w, err)
			return
		}
		input.EmployeeID = &parsed
	}

	entries, err := h.bookingUC.GetCalendar(ctx, orgCtx, input)
	if err != nil {
		util.SendError(ctx, w, err, bookingErrors)
		return
	}

	resp := getCalendarResponse{
		Calendar: make([]calendarEntryResponse, 0, len(entries)),
	}
	for _, e := range entries {
		resp.Calendar = append(resp.Calendar, calendarEntryResponse{
			AppointmentID:   e.AppointmentID,
			Status:          string(e.Status),
			StartAt:         e.StartAt,
			EndAt:           e.EndAt,
			Price:           e.Price.InexactFloat64(),
			DurationMinutes: e.DurationMinutes,
			CustomerComment: e.CustomerComment,
			EmployeeID:      e.EmployeeID,
			EmployeeName:    e.EmployeeName,
			CustomerID:      e.CustomerID,
			CustomerName:    e.CustomerName,
			CustomerPhone:   e.CustomerPhone,
			ServiceID:       e.ServiceID,
			ServiceName:     e.ServiceName,
			ServiceColor:    e.ServiceColor,
			LocationID:      e.LocationID,
			LocationName:    e.LocationName,
		})
	}

	coreHTTP.SendData(ctx, w, resp, http.StatusOK)
}
