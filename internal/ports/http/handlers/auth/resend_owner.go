package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/auth"
	coreHTTP "github.com/Rasikrr/core/http"
)

// Resend handles POST /api/v1/auth/register/resend.
//
// @Summary      Повторная отправка OTP
// @Description  Повторно отправляет OTP-код на телефон. Доступно только если прошёл retry_after интервал с предыдущей отправки.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      resendRequest    true  "Номер телефона"
// @Success      200   {object}  resendResponse
// @Failure      400   {object}  util.errorResponse
// @Failure      404   {object}  util.errorResponse  "Pending-запрос не найден"
// @Failure      429   {object}  util.errorResponse  "Слишком рано для повторной отправки"
// @Failure      500   {object}  util.errorResponse
// @Router       /api/v1/auth/register/resend [post]
func (h *Handler) resendOwner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req resendRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		util.SendBadRequest(ctx, w, err)
		return
	}

	out, err := h.registerOwnerUseCase.Resend(ctx, uc.ResendInput{
		Phone: req.Phone,
	})
	if err != nil {
		util.SendError(ctx, w, err, authErrors)
		return
	}

	coreHTTP.SendData(ctx, w, resendResponse{
		Message:    "code_sent",
		ExpiresIn:  out.ExpiresIn,
		RetryAfter: out.RetryAfter,
	}, http.StatusOK)
}
