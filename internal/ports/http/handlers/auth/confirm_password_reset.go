package auth

import (
	"net/http"

	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/auth"
	coreHTTP "github.com/Rasikrr/core/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
)

// confirmPasswordReset handles POST /api/v1/auth/password/reset/confirm.
//
// @Summary      Подтверждение сброса пароля (шаг 2/2)
// @Description  Проверяет one-time токен и устанавливает новый пароль. Все активные сессии инвалидируются.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      confirmResetRequest  true  "Токен и новый пароль"
// @Success      200   {object}  messageResponse
// @Failure      400   {object}  util.errorResponse
// @Failure      401   {object}  util.errorResponse  "Токен не найден или истёк"
// @Failure      500   {object}  util.errorResponse
// @Router       /api/v1/auth/password/reset/confirm [post]
func (h *Handler) confirmPasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req confirmResetRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		util.SendBadRequest(ctx, w, err)
		return
	}

	err := h.resetPasswordUseCase.ConfirmReset(ctx, uc.ConfirmResetInput{
		Token:       req.Token,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		util.SendError(ctx, w, err, authErrors)
		return
	}

	coreHTTP.SendData(ctx, w, &messageResponse{Message: "password_changed"}, http.StatusOK)
}
