// nolint: godot
package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/session"
	"github.com/Rasikrr/core/api"
	coreErr "github.com/Rasikrr/core/errors"
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
		api.SendError(w, errors.ErrSessionNotFound)
		return
	}

	valid, err := c.authService.ValidateRegisterToken(ctx, token)
	if err != nil {
		api.SendError(w, err)
		return
	}
	if !valid {
		api.SendError(w, coreErr.NewError("invalid or expired token", http.StatusUnauthorized))
		return
	}

	var req registerConfirmRequest
	if err = api.GetData(r, &req); err != nil {
		api.SendError(w, coreErr.ErrBadRequestBody.Wrap(err))
		return
	}

	tokens, err := c.authService.RegisterConfirm(ctx, req.Phone, req.Password)
	if err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, api.NewSuccessResponse(loginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}), http.StatusOK)
}
