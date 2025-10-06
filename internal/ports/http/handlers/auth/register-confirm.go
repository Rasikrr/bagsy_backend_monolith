// nolint: godot
package auth

import (
	"errors"
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/cookies"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/session"
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
// @Success 200 {object} api.SuccessResponse{data=loginResponse} "Регистрация успешно подтверждена"
// @Failure 400 {object} api.ErrorResponse "Неверные данные"
// @Failure 401 {object} api.ErrorResponse "Неверный или просроченный токен"
// @Failure 500 {object} api.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/register/confirm [post]
func (c *Controller) registerConfirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token, err := session.GetAuthHeader(r)
	if err != nil {
		api.SendError(w, err)
		return
	}

	valid, err := c.authService.ValidateRegisterToken(ctx, token)
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
	cookies.SetAuthTokens(w, tokens)
	api.SendData(w, api.NewSuccessResponse(loginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}), http.StatusOK)
}
