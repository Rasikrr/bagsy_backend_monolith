package pointcategories

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// getCategories godoc
// @Summary Получить список категорий точек
// @Description Возвращает список всех категорий точек
// @Tags point-categories
// @Accept json
// @Produce json
// @Success 200 {object} getCategoriesResponse "Список категорий точек"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/point-categories [get]
func (c *Controller) getCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := c.pointCategoriesService.GetAll(ctx)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, newGetCategoriesResponse(categories), http.StatusOK)
}
