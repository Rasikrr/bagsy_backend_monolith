package location

import (
	"net/http"
	"time"

	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
)

// getBySlug handles GET /api/v1/locations/slug/{slug}.
//
// @Summary      Получение локации по slug (публичный)
// @Description  Возвращает публично доступную локацию по её slug + расписание на 7 дней.
// @Tags         locations
// @Produce      json
// @Param        slug  path      string  true  "Slug локации"
// @Success      200  {object}  locationResponse
// @Failure      404  {object}  httputil.errorResponse
// @Failure      500  {object}  httputil.errorResponse
// @Router       /api/v1/locations/slug/{slug} [get]
func (h *Handler) getBySlug(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := chi.URLParam(r, "slug")

	loc, err := h.locationUseCase.GetBySlug(ctx, slug)
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
