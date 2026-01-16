package bagsies

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// resendCode godoc
// @Summary Повторная отправка кода подтверждения
// @Description Повторно отправляет код подтверждения брони клиенту по SMS/WhatsApp. Используется если клиент не получил код или код истек.
// @Tags bagsies
// @Accept json
// @Produce json
// @Param request body resentCodeRequest true "Данные для повторной отправки кода"
// @Success 200 {object} response.EmptySuccessResponse "Код успешно отправлен повторно"
// @Failure 400 {object} errors.ErrorResponse "Неверный формат запроса или валидация не пройдена"
// @Failure 404 {object} errors.ErrorResponse "Бронь не найдена"
// @Failure 429 {object} errors.ErrorResponse "Превышен лимит запросов на отправку кода"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/bagsies/resend [post]
func (c *Controller) resendCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req resentCodeRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	err := c.bagsiesService.ResendConfirmationCode(ctx, req.BagsyID)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, response.NewEmptySuccessResponse("confirmation code resent successfully"), http.StatusOK)
}
