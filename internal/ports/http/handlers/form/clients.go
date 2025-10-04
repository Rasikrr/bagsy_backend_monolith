// nolint: godot
package form

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/Rasikrr/core/api"
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
	var req clientFormRequest
	if err := api.GetData(r, &req); err != nil {
		api.SendError(w, err)
		return
	}
	if err := req.validate(); err != nil {
		api.SendError(w, err)
		return
	}
	ctx := r.Context()
	role, _ := enum.RoleString(req.Role)
	err := c.formsService.CreateClient(ctx, req.FirstName, req.LastName, req.Phone, req.Description, role)
	if err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, api.NewEmptySuccessResponse(), http.StatusOK)
}
