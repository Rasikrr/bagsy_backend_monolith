package location

import (
	"net/http"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// getByID handles GET /api/v1/locations/{id}.
//
// @Summary      Получение локации по ID
// @Description  Возвращает детали конкретной локации.
// @Tags         locations
// @Produce      json
// @Param        id   path      string  true  "UUID локации"
// @Success      200  {object}  locationResponse
// @Failure      400  {object}  httputil.errorResponse
// @Failure      401  {object}  httputil.errorResponse
// @Failure      403  {object}  httputil.errorResponse
// @Failure      404  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/locations/{id} [get]
func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	loc, err := h.locationUseCase.GetByID(ctx, orgCtx, id)
	if err != nil {
		httputil.SendError(ctx, w, err, locationErrors)
		return
	}

	now := time.Now()
	slots, err := h.scheduleRepo.GetLocationSlots(ctx, loc.ID, now, now.AddDate(0, 0, 7))
	if err != nil {
		httputil.SendError(ctx, w, err, locationErrors)
		return
	}

	coreHTTP.SendData(ctx, w, toLocationResponse(loc, slots), http.StatusOK)
}
