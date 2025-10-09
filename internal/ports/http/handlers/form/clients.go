// nolint: godot
package form

import (
	"net/http"

	"github.com/Rasikrr/core/api"
	coreErr "github.com/Rasikrr/core/errors"
)

// @Summary Отправка формы для соотрдничества
// @Description Создает новую заявку на сотрудничество
// @Tags forms
// @Accept json
// @Produce json
// @Param request body clientFormRequest true "Форма для создания заявки"
// @Success 200 {object} api.EmptySuccessResponse "Успешное создание заявки"
// @Failure 500 {object} api.ErrorResponse "ошибка"
// @Router /api/v1/forms [post]
func (c *Controller) createClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req clientFormRequest
	if err := api.GetData(r, &req); err != nil {
		api.SendError(w, coreErr.ErrBadRequestBody.Wrap(err))
		return
	}
	if err := req.validate(); err != nil {
		api.SendError(w, coreErr.ErrBadRequestBody.Wrap(err))
		return
	}
	err := c.formsService.CreateClient(ctx, req.FirstName, req.LastName, req.Phone, req.Description, req.Role)
	if err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, api.NewEmptySuccessResponse("application created"), http.StatusOK)
}
