// nolint: godot
package networks

import (
	"net/http"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
	"github.com/go-chi/chi/v5"
)

// getNetwork godoc
// @Summary Получение информации о сети по коду
// @Description Возвращает детальную информацию о сети: название, описание, даты создания/обновления и автора создания
// @Tags networks
// @Accept json
// @Produce json
// @Param code path string true "Код сети (slug)"
// @Success 200 {object} networkResponse "Информация о сети"
// @Failure 400 {object} errors.ErrorResponse "Неверный формат запроса"
// @Failure 404 {object} errors.ErrorResponse "Сеть не найдена"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/networks/{code} [get]
func (c *Controller) getNetwork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	if code == "" {
		errors.HandleError(ctx, w, domainErr.NewInvalidInputError("code parameter is required", nil))
		return
	}

	network, err := c.networksService.GetByCode(ctx, code)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, toNetworkResponse(network), http.StatusOK)
}
