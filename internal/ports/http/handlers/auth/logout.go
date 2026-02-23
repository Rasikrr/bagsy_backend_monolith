package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
)

// logout handles POST /api/v1/auth/logout.
//
// @Summary      Выход из системы
// @Description  Инвалидирует refresh-токен. Access-токен продолжает действовать до истечения TTL.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      logoutRequest    true  "Refresh-токен для инвалидации"
// @Success      204   "Токен успешно инвалидирован"
// @Failure      400   {object}  httputil.errorResponse
// @Failure      500   {object}  httputil.errorResponse
// @Router       /api/v1/auth/logout [post]
func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req logoutRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	if err := h.authUseCase.Logout(ctx, req.RefreshToken); err != nil {
		httputil.SendError(ctx, w, err, authErrors)
		return
	}

	coreHTTP.SendData(ctx, w, nil, http.StatusNoContent)
}
