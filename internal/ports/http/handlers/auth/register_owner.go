package auth

import (
	"net/http"

	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/auth"
	coreHTTP "github.com/Rasikrr/core/http"
)

// Register handles POST /api/v1/auth/register.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req registerRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
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
