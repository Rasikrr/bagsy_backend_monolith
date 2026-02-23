package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/auth"
	coreHTTP "github.com/Rasikrr/core/http"
)

// Register handles POST /api/v1/auth/register.
//
// @Summary      Регистрация владельца
// @Description  Начинает процесс регистрации владельца организации. Создаёт pending-запрос и отправляет OTP-код на указанный номер телефона.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      registerRequest   true  "Данные регистрации"
// @Success      200   {object}  registerResponse
// @Failure      400   {object}  httputil.errorResponse
// @Failure      409   {object}  httputil.errorResponse  "Пользователь с таким номером уже существует"
// @Failure      500   {object}  httputil.errorResponse
// @Router       /api/v1/auth/register [post]
func (h *Handler) registerOwner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req registerRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}
	out, err := h.registerOwnerUseCase.Register(ctx, uc.RegisterInput{
		Phone:     req.Phone,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
		PlanCode:  req.PlanCode,
	})
	if err != nil {
		httputil.SendError(ctx, w, err, authErrors)
		return
	}

	coreHTTP.SendData(ctx, w, registerResponse{
		Message:    "code_sent",
		Phone:      out.Phone,
		ExpiresIn:  out.ExpiresIn,
		RetryAfter: out.RetryAfter,
	}, http.StatusOK)
}
