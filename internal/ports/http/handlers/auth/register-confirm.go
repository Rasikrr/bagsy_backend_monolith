// nolint: godot
package auth

import (
	"errors"
	"net/http"

	"github.com/Rasikrr/core/api"
)

// RegisterConfirm godoc
// @Summary Подтверждение регистрации
// @Description Подтверждает регистрацию пользователя и устанавливает пароль
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Токен подтверждения регистрации"
// @Param request body registerConfirmRequest true "Данные для подтверждения"
// @Success 200 {object} map[string]string "Регистрация успешно подтверждена"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 401 {object} map[string]string "Неверный или просроченный токен"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/v1/auth/register/confirm [post]
func (c *Controller) registerConfirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.Header.Get("Authorization")
	if token == "" {
		api.SendError(w, errors.New("token is required"))
		return
	}

	valid, err := c.authService.ValidateRegistrationToken(ctx, token)
	if err != nil {
		api.SendError(w, err)
		return
	}
	if !valid {
		api.SendError(w, errors.New("invalid token"))
		return
	}

	var req registerConfirmRequest
	if getDataErr := api.GetData(r, &req); getDataErr != nil {
		api.SendError(w, getDataErr)
		return
	}

	tokens, err := c.authService.RegisterConfirm(ctx, req.Phone, req.Password)
	if err != nil {
		api.SendError(w, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
	})
	api.SendData(w, api.NewEmptySuccessResponse(), http.StatusOK)
}
