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

// getCalendar handles GET /api/v1/bookings/calendar.
//
// @Summary      Получение календаря записей
// @Description  Возвращает записи за указанный период с информацией о сотруднике, клиенте, услуге и локации.
// @Description  Staff — только свои записи. Manager — записи своей локации. Owner — все записи организации.
// @Description  Максимальный диапазон — 35 дней.
// @Tags         booking
// @Accept       json
// @Produce      json
// @Param        from               query     string  true   "Дата начала (YYYY-MM-DD)"
// @Param        to                 query     string  true   "Дата окончания (YYYY-MM-DD)"
// @Param        location_id        query     string  false  "UUID локации (только для Owner)"
// @Param        employee_id        query     string  false  "UUID сотрудника (для Manager и Owner)"
// @Param        include_cancelled  query     bool    false  "Включить отменённые записи (default: false)"
// @Success      200  {object}  getCalendarResponse
// @Failure      400  {object}  util.errorResponse  "Неверные параметры или диапазон > 35 дней"
// @Failure      401  {object}  util.errorResponse  "Требуется авторизация"
// @Failure      403  {object}  util.errorResponse  "Подписка приостановлена"
// @Failure      500  {object}  util.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/bookings/calendar [get]
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
		var parsed uuid.UUID
		parsed, err = uuid.Parse(locID)
		if err != nil {
			util.SendBadRequest(ctx, w, err)
			return
		}
		input.LocationID = &parsed
	}

	if empID := q.Get("employee_id"); empID != "" {
		var parsed uuid.UUID
		parsed, err = uuid.Parse(empID)
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
			Price:           e.Price.Amount().InexactFloat64(),
			DurationMinutes: e.DurationMinutes.Minutes(),
			CustomerComment: e.CustomerComment,
			EmployeeID:      e.EmployeeID,
			EmployeeName:    e.EmployeeName,
			CustomerID:      e.CustomerID,
			CustomerName:    e.CustomerName,
			CustomerPhone:   e.CustomerPhone.String(),
			ServiceID:       e.ServiceID,
			ServiceName:     e.ServiceName,
			ServiceColor:    string(e.ServiceColor),
			LocationID:      e.LocationID,
			LocationName:    e.LocationName,
		})
	}

	coreHTTP.SendData(ctx, w, resp, http.StatusOK)
}
