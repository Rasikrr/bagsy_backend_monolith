package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// Refresh godoc
// @Summary Обновление токенов авторизации
// @Description Обновляет токены авторизации
// @Tags auth
// @Accept json
// @Produce json
// @Param request body refreshTokensRequest true "Токен авторизации"
// @Success 200 {object} refreshTokensResponse "Токены авторизации успешно обновлены"
// @Failure 400 {object} errors.ErrorResponse "Неверный токен"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/refresh [post]
func (c *Controller) refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req refreshTokensRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	access, refresh, err := c.authService.Refresh(ctx, req.RefreshToken)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	response.SendData(ctx, w, refreshTokensResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, http.StatusOK)
}
