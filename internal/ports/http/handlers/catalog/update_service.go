package catalog

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/catalog"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// updateService handles PUT /api/v1/services/{id}.
//
// @Summary      Обновление услуги
// @Description  Обновляет данные услуги. Все поля опциональны. Owner — любую в org, Manager — в своей локации.
// @Tags         catalog
// @Accept       json
// @Param        id    path  string               true  "ID услуги"
// @Param        body  body  updateServiceRequest  true  "Данные для обновления"
// @Success      204
// @Failure      400  {object}  httputil.errorResponse
// @Failure      403  {object}  httputil.errorResponse
// @Failure      404  {object}  httputil.errorResponse
// @Failure      410  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/services/{id} [put]
func (h *Handler) updateService(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	serviceID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	var req updateServiceRequest
	if err = coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	input := uc.UpdateServiceInput{
		ID:              serviceID,
		Name:            req.Name,
		Description:     req.Description,
		Color:           req.Color,
		DurationMinutes: req.DurationMinutes,
		SortOrder:       req.SortOrder,
	}

	if err = h.catalogUseCase.UpdateService(ctx, orgCtx, input); err != nil {
		httputil.SendError(ctx, w, err, catalogErrors)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
