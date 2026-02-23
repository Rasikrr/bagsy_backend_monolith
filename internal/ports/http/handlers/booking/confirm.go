package booking

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// confirm handles POST /api/v1/bookings/{id}/confirm.
//
// @Summary      Подтверждение записи
// @Description  Подтверждает запись с помощью OTP-кода, присланного в SMS/WhatsApp.
// @Tags         booking
// @Accept       json
// @Produce      json
// @Param        id    path      string          true  "ID записи"
// @Param        body  body      confirmRequest  true  "Код подтверждения"
// @Success      204   "Запись подтверждена"
// @Failure      400   {object}  util.errorResponse  "Неверный код или формат"
// @Failure      404   {object}  util.errorResponse  "Запись не найдена"
// @Failure      500   {object}  util.errorResponse
// @Router       /api/v1/bookings/{id}/confirm [post]
func (h *Handler) confirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		util.SendBadRequest(ctx, w, err)
		return
	}

	var req confirmRequest
	if err = coreHTTP.GetData(r, &req); err != nil {
		util.SendBadRequest(ctx, w, err)
		return
	}

	if err = h.bookingUC.Confirm(ctx, id, req.Code); err != nil {
		util.SendError(ctx, w, err, bookingErrors)
		return
	}

	coreHTTP.SendData(ctx, w, nil, http.StatusNoContent)
}
