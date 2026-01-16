// nolint
package points

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/dto"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// createPoint godoc
// @Summary Создать точку обслуживания
// @Description Создаёт новую точку обслуживания с указанными параметрами. Опционально можно прикрепить фото (до 10). Создавать могут только NetManager/SelfOwner
// @Tags points
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body createPointRequest true "Данные для создания точки (photo_ids опционально)"
// @Success 201 {object} dto.PointResponse "Точка успешно создана"
// @Failure 400 {object} errors.ErrorResponse "Неверные параметры запроса или некорректный формат photo_id"
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
	cmd, err := req.toCommand()
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	point, err := c.pointsService.Create(ctx, cmd)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	response.SendData(ctx, w, dto.ToPointResponse(point), http.StatusCreated)
}
