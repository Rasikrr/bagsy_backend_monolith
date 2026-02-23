package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
)

// Login handles POST /api/v1/auth/login.
//
// @Summary      Вход сотрудника
// @Description  Аутентификация сотрудника по номеру телефона и паролю. Возвращает пару JWT-токенов (access + refresh).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      loginRequest     true  "Телефон и пароль"
// @Success      200   {object}  loginResponse
// @Failure      400   {object}  httputil.errorResponse
// @Failure      401   {object}  httputil.errorResponse  "Неверный телефон или пароль"
// @Failure      500   {object}  httputil.errorResponse
// @Router       /api/v1/auth/login [post]
func (h *Handler) loginEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req loginRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}
	out, err := h.authUseCase.LoginEmployee(ctx, req.Phone, req.Password)
	if err != nil {
		httputil.SendError(ctx, w, err, authErrors)
		return
	}
	coreHTTP.SendData(ctx, w, &loginResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	}, http.StatusOK)
}
