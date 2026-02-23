package booking

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// cancel handles POST /api/v1/bookings/{id}/cancel.
//
// @Summary      Отмена записи
// @Description  Отменяет существующую запись. Доступно только сотрудникам организации.
// @Tags         booking
// @Accept       json
// @Produce      json
// @Param        id    path      string         true  "ID записи"
// @Param        body  body      cancelRequest  true  "Причина отмены"
// @Success      204   "Запись отменена"
// @Failure      400   {object}  httputil.errorResponse
// @Failure      403   {object}  httputil.errorResponse  "Доступ запрещен"
// @Failure      404   {object}  httputil.errorResponse  "Запись не найдена"
// @Failure      500   {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/bookings/{id}/cancel [post]
func (h *Handler) cancel(w http.ResponseWriter, r *http.Request) {
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

	var req cancelRequest
	if err = coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	if err = h.bookingUC.Cancel(ctx, orgCtx, id, req.Reason); err != nil {
		httputil.SendError(ctx, w, err, bookingErrors)
		return
	}

	coreHTTP.SendData(ctx, w, nil, http.StatusNoContent)
}
