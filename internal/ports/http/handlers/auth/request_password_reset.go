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
// @Failure      400   {object}  util.errorResponse
// @Failure      403   {object}  util.errorResponse  "Аккаунт неактивен"
// @Failure      404   {object}  util.errorResponse  "Сотрудник не найден"
// @Failure      500   {object}  util.errorResponse
// @Router       /api/v1/auth/password/reset [post]
func (h *Handler) requestPasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req requestResetRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		util.SendBadRequest(ctx, w, err)
		return
	}

	err := h.resetPasswordUseCase.RequestReset(ctx, uc.RequestResetInput{
		Phone: req.Phone,
	})
	if err != nil {
		util.SendError(ctx, w, err, authErrors)
		return
	}

	coreHTTP.SendData(ctx, w, &messageResponse{Message: "reset_link_sent"}, http.StatusOK)
}
