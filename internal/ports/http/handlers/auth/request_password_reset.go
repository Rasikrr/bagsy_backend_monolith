package auth

import (
	"net/http"

	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/auth"
	coreHTTP "github.com/Rasikrr/core/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
)

// requestPasswordReset handles POST /api/v1/auth/password/reset.
//
// @Summary      Запрос на сброс пароля (шаг 1/2)
// @Description  Отправляет ссылку для сброса пароля на номер телефона сотрудника через WhatsApp (с fallback на SMS).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      requestResetRequest  true  "Номер телефона"
// @Success      200   {object}  messageResponse
// @Failure      400   {object}  httputil.errorResponse
// @Failure      403   {object}  httputil.errorResponse  "Аккаунт неактивен"
// @Failure      404   {object}  httputil.errorResponse  "Сотрудник не найден"
// @Failure      500   {object}  httputil.errorResponse
// @Router       /api/v1/auth/password/reset [post]
func (h *Handler) requestPasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req requestResetRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	err := h.resetPasswordUseCase.RequestReset(ctx, uc.RequestResetInput{
		Phone: req.Phone,
	})
	if err != nil {
		httputil.SendError(ctx, w, err, authErrors)
		return
	}

	coreHTTP.SendData(ctx, w, &messageResponse{Message: "reset_link_sent"}, http.StatusOK)
}
