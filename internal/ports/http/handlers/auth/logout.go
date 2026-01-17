package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// logout godoc
// @Summary Выход из системы
// @Description Выполняет выход пользователя из системы и инвалидирует refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body refreshTokensRequest true "Refresh token для выхода"
// @Success 200 "Успешный выход"
// @Failure 400 {object} errors.ErrorResponse "Неверные данные"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/logout [post]
func (c *Controller) logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req refreshTokensRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	if err := c.authService.Logout(ctx, req.RefreshToken); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, nil, http.StatusOK)
}
