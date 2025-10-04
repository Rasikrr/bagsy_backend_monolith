package auth

import (
	"errors"
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/cookies"
	"github.com/Rasikrr/core/api"
)

// Refresh godoc
// @Summary Обновление токенов авторизации
// @Description Обновляет токены авторизации
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Токен авторизации"
// @Success 200 {object} map[string]string "Токены авторизации успешно обновлены"
// @Failure 400 {object} map[string]string "Неверный токен"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/v1/auth/refresh [post]
// nolint: godot
func (c *Controller) refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := cookies.GetRefreshToken(r)
	if token == "" {
		api.SendError(w, errors.New("token is required"))
		return
	}
	tokens, err := c.authService.RefreshTokens(ctx, token)
	if err != nil {
		api.SendError(w, err)
		return
	}
	cookies.SetAuthTokens(w, tokens)
	api.SendData(w, api.NewEmptySuccessResponse(), http.StatusOK)
}
