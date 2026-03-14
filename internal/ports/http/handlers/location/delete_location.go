package location

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// deleteLocation handles DELETE /api/v1/locations/{id}.
//
// @Summary      Удаление локации
// @Description  Soft-delete локации. Только Owner.
// @Tags         locations
// @Param        id  path  string  true  "ID локации"
// @Success      204
// @Failure      403  {object}  httputil.errorResponse
// @Failure      404  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/locations/{id} [delete]
func (h *Handler) deleteLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	locationID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	if err = h.locationUseCase.DeleteLocation(ctx, orgCtx, locationID); err != nil {
		httputil.SendError(ctx, w, err, locationErrors)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
