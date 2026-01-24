package masterservices

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// createMasterService godoc
// @Summary Создать связь мастер-услуга
// @Description Создаёт связь между мастером и услугой с указанной ценой. Staff создаёт для себя, Manager — для сотрудников своей точки, NetManager/SelfOwner — для сотрудников своей сети.
// @Tags master-services
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body createMasterServiceRequest true "Данные для создания связи мастер-услуга"
// @Success 201 {object} createMasterServiceResponse "Связь успешно создана"
// @Failure 400 {object} errors.ErrorResponse "Неверные параметры запроса"
// @Failure 401 {object} errors.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} errors.ErrorResponse "Недостаточно прав"
// @Failure 409 {object} errors.ErrorResponse "Связь мастер-услуга уже существует"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/master-services [post]
func (c *Controller) createMasterService(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req createMasterServiceRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	ms, err := c.masterServicesService.Create(ctx, req.toCommand())
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, newCreateMasterServiceResponse(ms), http.StatusCreated)
}
