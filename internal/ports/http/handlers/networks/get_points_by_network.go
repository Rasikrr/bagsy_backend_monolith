// nolint: godot
package networks

import (
	"net/http"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/dto"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
	"github.com/go-chi/chi/v5"
)

// getPointsByNetwork godoc
// @Summary Получение списка точек сети
// @Description Возвращает список всех точек (заведений), принадлежащих указанной сети. Каждая точка содержит информацию о названии, адресе, расписании работы и других характеристиках.
// @Tags networks
// @Accept json
// @Produce json
// @Param code path string true "Код сети (slug)"
// @Success 200 {object} dto.PointsResponse "Список точек сети"
// @Failure 400 {object} errors.ErrorResponse "Неверный формат запроса"
// @Failure 404 {object} errors.ErrorResponse "Сеть не найдена"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/networks/{code}/points [get]
func (c *Controller) getPointsByNetwork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	if code == "" {
		errors.HandleError(ctx, w, domainErr.NewInvalidInputError("code parameter is required", nil))
		return
	}
	points, err := c.pointsService.GetByNetworkCode(ctx, code)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	response.SendData(ctx, w, dto.ToPointsResponse(points), http.StatusOK)
}
