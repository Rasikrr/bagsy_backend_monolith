package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/auth"
	coreHTTP "github.com/Rasikrr/core/http"
)

// Verify handles POST /api/v1/auth/register/verify.
//
// @Summary      Подтверждение регистрации
// @Description  Подтверждает регистрацию владельца OTP-кодом. При успехе создаёт организацию, сотрудника, подписку и возвращает JWT-токены.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      verifyRequest    true  "Телефон и OTP-код"
// @Success      200   {object}  tokensResponse
// @Failure      400   {object}  util.errorResponse  "Неверный код или истёк срок действия"
// @Failure      404   {object}  util.errorResponse  "Pending-запрос не найден"
// @Failure      500   {object}  util.errorResponse
// @Router       /api/v1/auth/register/verify [post]
func (h *Handler) verifyOwner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req verifyRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		util.SendBadRequest(ctx, w, err)
		return
	}

	out, err := h.registerOwnerUseCase.VerifyRegistration(ctx, uc.VerifyInput{
		Phone: req.Phone,
		Code:  req.Code,
	})
	if err != nil {
		util.SendError(ctx, w, err, authErrors)
		return
	}

	coreHTTP.SendData(ctx, w, tokensResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	}, http.StatusOK)
}
