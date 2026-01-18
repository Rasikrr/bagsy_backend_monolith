// nolint: godot
package calendar

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// getCalendar godoc
// @Summary Получение календаря записей
// @Description Возвращает записи за указанный период.
// @Description Для Staff - только свои записи.
// @Description Для Manager+ - записи всей точки (опционально с фильтром по мастеру).
// @Tags calendar
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param from query string true "Дата начала в формате YYYY-MM-DD"
// @Param to query string true "Дата окончания в формате YYYY-MM-DD"
// @Param point_code query string false "Код точки для фильтрации (только для Manager+)"
// @Param master_phone query string false "Телефон мастера для фильтрации (только для Manager+)"
// @Success 200 {object} calendarResponseDTO
// @Failure 400 {object} errors.ErrorResponse "Неверные параметры запроса"
// @Failure 401 {object} errors.ErrorResponse "Требуется авторизация"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/calendar [get]
func (c *Controller) getCalendar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req getCalendarRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	elements, err := c.calendarService.GetCalendar(ctx, req.toQuery())
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, newCalendarResponse(elements), http.StatusOK)
}
