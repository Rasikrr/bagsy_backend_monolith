package location

import (
	"net/http"

	domainLoc "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
)

// getCategories handles GET /api/v1/location-categories.
//
// @Summary      Список категорий локаций
// @Description  Возвращает все доступные категории локаций, отсортированные по sort_order.
// @Tags         locations
// @Produce      json
// @Success      200  {object}  getCategoriesResponse
// @Failure      500  {object}  httputil.errorResponse
// @Router       /api/v1/locations/categories [get]
func (h *Handler) getCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categories, err := h.locationUseCase.GetCategories(ctx)
	if err != nil {
		httputil.SendError(ctx, w, err, locationErrors)
		return
	}
	coreHTTP.SendData(ctx, w, toGetCategoriesResponse(categories), http.StatusOK)
}

func toGetCategoriesResponse(cats []*domainLoc.Category) getCategoriesResponse {
	items := make([]locationCategoryResponse, 0, len(cats))
	for _, c := range cats {
		items = append(items, locationCategoryResponse{
			ID:        c.ID.String(),
			Slug:      c.Slug.String(),
			Name:      c.Name,
			SortOrder: c.SortOrder,
		})
	}
	return getCategoriesResponse{Categories: items}
}
