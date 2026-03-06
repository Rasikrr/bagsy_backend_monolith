package catalog

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/catalog"
	coreHTTP "github.com/Rasikrr/core/http"
)

var errLocationCategoryIDRequired = errors.New("location_category_id is required")

// getServiceCategories handles GET /api/v1/service-categories.
//
// @Summary      Получение категорий услуг
// @Description  Возвращает дерево категорий услуг для указанного типа бизнеса (location_category_id).
// @Tags         catalog
// @Produce      json
// @Param        location_category_id  query     string  true  "ID категории локации"
// @Success      200  {object}  getServiceCategoriesResponse
// @Failure      400  {object}  httputil.errorResponse
// @Failure      401  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/service-categories [get]
func (h *Handler) getServiceCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	locationCategoryIDStr := r.URL.Query().Get("location_category_id")
	if locationCategoryIDStr == "" {
		httputil.SendBadRequest(ctx, w, errLocationCategoryIDRequired)
		return
	}

	locationCategoryID, err := uuid.Parse(locationCategoryIDStr)
	if err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	trees, err := h.catalogUseCase.GetServiceCategories(ctx, locationCategoryID)
	if err != nil {
		httputil.SendError(ctx, w, err, catalogErrors)
		return
	}

	coreHTTP.SendData(ctx, w, toGetServiceCategoriesResponse(trees), http.StatusOK)
}

func toGetServiceCategoriesResponse(trees []uc.ServiceCategoryTree) getServiceCategoriesResponse {
	categories := make([]serviceCategoryResponse, 0, len(trees))
	for _, t := range trees {
		categories = append(categories, toServiceCategoryResponse(t))
	}
	return getServiceCategoriesResponse{Categories: categories}
}

func toServiceCategoryResponse(t uc.ServiceCategoryTree) serviceCategoryResponse {
	children := make([]serviceCategoryResponse, 0, len(t.Children))
	for _, c := range t.Children {
		children = append(children, toServiceCategoryResponse(c))
	}
	return serviceCategoryResponse{
		ID:        t.ID.String(),
		Name:      t.Name,
		SortOrder: t.SortOrder,
		Children:  children,
	}
}
