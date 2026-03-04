package catalog

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/catalog"
	coreHTTP "github.com/Rasikrr/core/http"
)

// createService handles POST /api/v1/services.
//
// @Summary      Создание услуги
// @Description  Создаёт новую услугу в локации. Owner — в любой, Manager — в своей.
// @Tags         catalog
// @Accept       json
// @Produce      json
// @Param        body  body      createServiceRequest  true  "Данные услуги"
// @Success      201   {object}  createServiceResponse
// @Failure      400   {object}  httputil.errorResponse
// @Failure      403   {object}  httputil.errorResponse
// @Failure      404   {object}  httputil.errorResponse
// @Failure      500   {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/services [post]
func (h *Handler) createService(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	var req createServiceRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	input := uc.CreateServiceInput{
		LocationID:      req.LocationID,
		CategoryID:      req.CategoryID,
		Name:            req.Name,
		Description:     req.Description,
		Color:           req.Color,
		DurationMinutes: req.DurationMinutes,
	}

	out, err := h.catalogUseCase.CreateService(ctx, orgCtx, input)
	if err != nil {
		httputil.SendError(ctx, w, err, catalogErrors)
		return
	}

	coreHTTP.SendData(ctx, w, createServiceResponse{
		ID: out.ID.String(),
	}, http.StatusCreated)
}
