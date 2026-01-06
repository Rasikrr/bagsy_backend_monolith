// nolint: godot
package points

import (
	"net/http"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/dto"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
	"github.com/go-chi/chi/v5"
)

// getPoint godoc
// @Summary Получить точку обслуживания
// @Description Возвращает информацию о точке обслуживания по её коду
// @Tags points
// @Accept json
// @Produce json
// @Param code path string true "Код точки обслуживания"
// @Success 200 {object} dto.PointResponse "Информация о точке"
// @Failure 400 {object} errors.ErrorResponse "Неверные параметры запроса"
// @Failure 401 {object} errors.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} errors.ErrorResponse "Недостаточно прав для доступа"
// @Failure 404 {object} errors.ErrorResponse "Точка не найдена"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/points/{code} [get]
func (c *Controller) getPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")
	if code == "" {
		errors.HandleError(ctx, w, domainErr.NewInvalidInputError("code parameter is required", nil))
		return
	}

	point, err := c.pointsService.GetByCode(ctx, code)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, dto.ToPointResponse(point), http.StatusOK)
}
