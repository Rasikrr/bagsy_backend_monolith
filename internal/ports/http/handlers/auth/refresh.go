package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/pkg/session"
	"github.com/Rasikrr/core/api"
)

// Refresh godoc
// @Summary Обновление токенов авторизации
// @Description Обновляет токены авторизации
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Токен авторизации"
// @Success 200 {object} api.EmptySuccessResponse "Токены авторизации успешно обновлены"
// @Failure 400 {object} api.ErrorResponse "Неверный токен"
// @Failure 500 {object} api.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/refresh [post]
// nolint: godot
func (c *Controller) refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token, err := session.GetAuthHeader(r)
	if err != nil {
		api.SendError(w, err)
		return
	}
	tokens, err := c.authService.RefreshTokens(ctx, token)
	if err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, api.NewSuccessResponse(refreshTokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}), http.StatusOK)
}
