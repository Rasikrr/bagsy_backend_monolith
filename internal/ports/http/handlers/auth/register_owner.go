package auth

import (
	"net/http"

	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/auth"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/Rasikrr/core/log"
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
// @Failure      400   {object}  coreHTTP.ErrorResponse
// @Failure      409   {object}  coreHTTP.ErrorResponse  "Пользователь с таким номером уже существует"
// @Failure      500   {object}  coreHTTP.ErrorResponse
// @Router       /api/v1/auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req registerRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		log.Error(ctx, "get data error", log.Err(err))
		// TODO: error
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
		log.Error(ctx, "failed to register", log.Err(err))

		// TODO: error
		return
	}

	coreHTTP.SendData(ctx, w, registerResponse{
		Message:    "code_sent",
		Phone:      out.Phone,
		ExpiresIn:  out.ExpiresIn,
		RetryAfter: out.RetryAfter,
	}, http.StatusOK)

}
