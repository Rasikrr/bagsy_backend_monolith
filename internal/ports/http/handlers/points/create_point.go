// nolint: godot
package points

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// createPoint godoc
// @Summary Создать точку обслуживания
// @Description Создаёт новую точку обслуживания с указанными параметрами. Создавать могут только NetManager/SelfOwner
// @Tags points
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body createPointRequest true "Данные для создания точки"
// @Success 201 {object} pointCreateResponse "Точка успешно создана"
// @Failure 400 {object} errors.ErrorResponse "Неверные параметры запроса"
// @Failure 401 {object} errors.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} errors.ErrorResponse "Недостаточно прав для создания точки"
// @Failure 409 {object} errors.ErrorResponse "Точка с таким кодом уже существует"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/points [post]
func (c *Controller) createPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req createPointRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	point, err := req.toEntity()
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	if err = c.pointsService.Create(ctx, point); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, toPointCreateResponse(point.Code), http.StatusCreated)
}
