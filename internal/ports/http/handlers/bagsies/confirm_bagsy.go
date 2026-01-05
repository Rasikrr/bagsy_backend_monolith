package bagsies

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// confirmBagsy godoc
// @Summary Подтверждение брони
// @Description Подтверждает созданную бронь с помощью кода, отправленного клиенту по SMS/WhatsApp
// @Tags bagsies
// @Accept json
// @Produce json
// @Param request body confirmBagsyRequest true "Данные для подтверждения брони"
// @Success 201 {object} response.EmptySuccessResponse "Бронь успешно подтверждена"
// @Failure 400 {object} errors.ErrorResponse "Неверный формат запроса или валидация не пройдена"
// @Failure 404 {object} errors.ErrorResponse "Бронь не найдена"
// @Failure 401 {object} errors.ErrorResponse "Неверный или просроченный код подтверждения"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/bagsies/confirm [post]
func (c *Controller) confirmBagsy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req confirmBagsyRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	err := c.bagsiesService.Confirm(ctx, req.BagsyID, req.Code)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	response.SendData(ctx, w, response.NewEmptySuccessResponse("confirmed"), http.StatusCreated)
}
