package bagsies

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// createBagsyByMaster godoc
// @Summary Создание брони мастером
// @Description Создает бронь напрямую от имени мастера без подтверждения кода. Статус сразу created.
// @Tags bagsies
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body createBagsyRequest true "Данные для создания брони"
// @Success 201 {object} createBagsyResponse "Бронь успешно создана"
// @Failure 400 {object} errors.ErrorResponse "Неверный формат запроса или валидация не пройдена"
// @Failure 401 {object} errors.ErrorResponse "Не авторизован"
// @Failure 403 {object} errors.ErrorResponse "Нет прав доступа"
// @Failure 404 {object} errors.ErrorResponse "Услуга или мастер не найдены"
// @Failure 409 {object} errors.ErrorResponse "Выбранное время уже занято у данного мастера"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/bagsies/master [post]
func (c *Controller) createBagsyByMaster(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req createBagsyRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	bagsyID, err := c.bagsiesService.CreateByMaster(ctx, req.toDomain())
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	response.SendData(ctx, w, newCreateBagsyResponse(bagsyID), http.StatusCreated)
}
