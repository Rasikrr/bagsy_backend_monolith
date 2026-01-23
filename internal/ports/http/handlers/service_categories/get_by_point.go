package servicecategories

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
	"github.com/go-chi/chi/v5"
)

// getByPointCode godoc
// @Summary Получить категории услуг по коду точки
// @Description Возвращает список категорий услуг и их подкатегорий, доступных для указанной точки
// @Tags service-categories
// @Accept json
// @Produce json
// @Param point_code path string true "Код точки"
// @Success 200 {object} getByPointCodeResponse "Список категорий услуг с подкатегориями"
// @Failure 404 {object} errors.ErrorResponse "Точка не найдена"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/service-categories/{point_code} [get]
func (c *Controller) getByPointCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pointCode := chi.URLParam(r, "point_code")

	categories, err := c.serviceCategoriesService.GetByPointCode(ctx, pointCode)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, newGetByPointCodeResponse(categories), http.StatusOK)
}
