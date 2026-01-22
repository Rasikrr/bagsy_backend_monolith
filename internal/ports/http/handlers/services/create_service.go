package services

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/actor"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// createService godoc
// @Summary Создать услугу
// @Description Создаёт новую услугу для точки обслуживания
// @Tags services
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body createServiceRequest true "Данные для создания услуги"
// @Success 201 {object} serviceDTO "Услуга успешно создана"
// @Failure 400 {object} errors.ErrorResponse "Неверные параметры запроса"
// @Failure 401 {object} errors.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} errors.ErrorResponse "Недостаточно прав для создания услуги"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/services [post]
func (c *Controller) createService(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	act, err := actor.GetActor(ctx)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	var req createServiceRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	cmd, err := req.toCommand(act.Phone())
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	id, err := c.servicesService.Create(ctx, cmd)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, newCreateServiceResponse(id), http.StatusCreated)
}
