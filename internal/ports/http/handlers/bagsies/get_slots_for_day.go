package bagsies

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// @Summary Получить слоты на конкретный день
// @Description Возвращает список доступных временных слотов для записи на услугу на конкретную дату, сгруппированных по мастерам
// @Tags bagsies
// @Accept json
// @Produce json
// @Param request body getSlotsForDayRequest true "Параметры запроса слотов на день"
// @Success 200 {object} getSlotsForDayResponse "Слоты на выбранный день"
// @Failure 400 {object} errors.ErrorResponse "Неверные параметры запроса"
// @Failure 404 {object} errors.ErrorResponse "Точка, услуга или мастера не найдены"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/bagsies/slots/day [post]
func (c *Controller) getSlotsForDay(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req getSlotsForDayRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	cmd, err := req.toDomain()
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	slots, err := c.bagsiesService.GetAvailableSlots(ctx, cmd)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, newGetSlotsForDayResponse(slots, req.Date), http.StatusOK)
}
