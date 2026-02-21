package booking

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// resendOTP handles POST /api/v1/bookings/{id}/resend-otp.
//
// @Summary      Повторная отправка кода
// @Description  Генерирует и отправляет новый код подтверждения для существующей записи.
// @Tags         booking
// @Produce      json
// @Param        id    path      string  true  "ID записи"
// @Success      204   "Код отправлен"
// @Failure      400   {object}  util.errorResponse  "Запись уже подтверждена или отменена"
// @Failure      404   {object}  util.errorResponse  "Запись не найдена"
// @Failure      500   {object}  util.errorResponse
// @Router       /api/v1/bookings/{id}/resend-otp [post]
func (h *Handler) resendOTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		util.SendBadRequest(ctx, w, err)
		return
	}

	if err := h.bookingUC.ResendOTP(ctx, id); err != nil {
		util.SendError(ctx, w, err, bookingErrors)
		return
	}

	coreHTTP.SendData(ctx, w, nil, http.StatusNoContent)
}
