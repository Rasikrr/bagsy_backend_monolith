// nolint: godot
package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// registerStaffConfirm godoc
// @Summary Подтверждение регистрации работника
// @Description Подтверждает регистрацию пользователя(работника)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body registerConfirmRequest true "Данные для подтверждения"
// @Success 200 {object} registerConfirmResponse "Регистрация успешно подтверждена"
// @Failure 400 {object} errors.ErrorResponse "Неверные данные"
// @Failure 401 {object} errors.ErrorResponse "Неверный или просроченный токен"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/staff/register/confirm [post]
func (c *Controller) registerStaffConfirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req registerConfirmRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	access, refresh, err := c.authService.RegisterStaffConfirm(ctx, req.toDomain())
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	response.SendData(ctx, w, registerConfirmResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, http.StatusOK)
}
