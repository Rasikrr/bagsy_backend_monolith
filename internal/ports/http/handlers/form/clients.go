// nolint: godot
package form

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// @Summary Отправка формы для сотрудничества
// @Description Создает новую заявку на сотрудничество
// @Tags forms
// @Accept json
// @Produce json
// @Param request body clientFormRequest true "Форма для создания заявки"
// @Success 200 {object} response.EmptySuccessResponse "Успешное создание заявки"
// @Failure 500 {object} errors.ErrorResponse "ошибка"
// @Router /api/v1/forms [post]
func (c *Controller) createClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req clientFormRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	err := c.formsService.Create(ctx, req.toEntity())
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, response.NewEmptySuccessResponse("ok"), http.StatusCreated)
}
