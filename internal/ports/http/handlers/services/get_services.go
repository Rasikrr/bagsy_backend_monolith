// nolint: godot
package services

import (
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/core/log"
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// getServicesByPointCode godoc
// @Summary Получить список услуг по коду точки
// @Description Возвращает список услуг для указанной точки. По умолчанию возвращает все услуги.
// @Tags services
// @Accept json
// @Produce json
// @Param point_code path string true "Код точки"
// @Param is_active query boolean false "Фильтр по активности услуги"
// @Success 200 {object} getServicesResponse "Список услуг"
// @Failure 400 {object} errors.ErrorResponse "Неверные параметры запроса"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/services/{point_code} [get]
func (c *Controller) getServicesByPointCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "point_code")
	if code == "" {
		errors.HandleError(ctx, w, domainErr.NewInvalidInputError("code parameter is required", nil))
		return
	}

	var req getServicesRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	log.Info(ctx, "req", log.Any("req", req))

	services, err := c.servicesService.GetByPointCode(ctx, code, req.Active)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, newGetServicesResponse(services), http.StatusOK)
}
