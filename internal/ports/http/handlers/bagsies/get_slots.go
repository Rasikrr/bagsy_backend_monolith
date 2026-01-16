package bagsies

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// @Summary Получить доступные слоты для записи
// @Description Возвращает список доступных временных слотов для записи на услугу на точке за ближайшие 2 недели
// @Tags bagsies
// @Accept json
// @Produce json
// @Param request body getSlotsRequest true "Параметры запроса слотов"
// @Success 200 {object} getSlotsResponse "Слоты казик"
// @Failure 400 {object} errors.ErrorResponse "Неверные параметры запроса"
// @Failure 404 {object} errors.ErrorResponse "Точка, услуга или мастера не найдены"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/bagsies/slots [post]
func (c *Controller) getSlots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req getSlotsRequest
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

	response.SendData(ctx, w, newGetSlotsResponse(slots), http.StatusOK)
}
