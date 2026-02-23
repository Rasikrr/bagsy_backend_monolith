package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
)

// refreshTokens handles POST /api/v1/auth/refresh.
//
// @Summary      Обновление токенов
// @Description  Принимает refresh-токен, ротирует его и возвращает новую пару access + refresh токенов.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      refreshRequest   true  "Refresh-токен"
// @Success      200   {object}  tokensResponse
// @Failure      400   {object}  httputil.errorResponse
// @Failure      401   {object}  httputil.errorResponse  "Токен не найден или истёк"
// @Failure      500   {object}  httputil.errorResponse
// @Router       /api/v1/auth/refresh [post]
func (h *Handler) refreshTokens(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req refreshRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	out, err := h.authUseCase.RefreshTokens(ctx, req.RefreshToken)
	if err != nil {
		httputil.SendError(ctx, w, err, authErrors)
		return
	}

	coreHTTP.SendData(ctx, w, tokensResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	}, http.StatusOK)
}
