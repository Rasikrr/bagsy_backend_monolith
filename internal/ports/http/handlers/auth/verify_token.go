package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
)

// verifyActionToken handles GET /api/v1/auth/verify/{token}.
//
// @Summary      Проверить action-токен
// @Description  Возвращает метаданные токена (телефон, назначение, опциональные org/location ID) для отображения на фронтенде. Используется для приглашений и сброса пароля. Не требует авторизации.
// @Tags         auth
// @Produce      json
// @Param        token  path      string  true  "Action-токен"
// @Success      200    {object}  verifyActionTokenResponse
// @Failure      400    {object}  util.errorResponse  "Пустой токен"
// @Failure      404    {object}  util.errorResponse  "Токен не найден"
// @Failure      500    {object}  util.errorResponse
// @Router       /api/v1/auth/verify/{token} [get]
func (h *Handler) verifyActionToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := chi.URLParam(r, "token")
	if token == "" {
		util.SendBadRequest(ctx, w, nil)
		return
	}

	actionToken, err := h.authUseCase.VerifyActionToken(ctx, token)
	if err != nil {
		util.SendError(ctx, w, err, authErrors)
		return
	}

	resp := verifyActionTokenResponse{
		Phone:   actionToken.Phone.String(),
		Purpose: actionToken.Purpose.String(),
	}
	if actionToken.LocationID != nil {
		id := actionToken.LocationID.String()
		resp.LocationID = &id
	}
	if actionToken.OrganizationID != nil {
		id := actionToken.OrganizationID.String()
		resp.OrganizationID = &id
	}

	coreHTTP.SendData(ctx, w, resp, http.StatusOK)
}
