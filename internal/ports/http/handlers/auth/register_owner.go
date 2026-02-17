package auth

import (
	"errors"
	"net/http"

	authDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
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
		mapRegisterError(w, err)
		return
	}

	coreHTTP.SendData(w, registerResponse{
		Message:    "code_sent",
		Phone:      out.Phone,
		ExpiresIn:  out.ExpiresIn,
		RetryAfter: out.RetryAfter,
	}, http.StatusOK)
}

//// ── Error mapping ───────────────────────────────────────────────
//
//func mapRegisterError(w http.ResponseWriter, err error) {
//	switch {
//	case errors.Is(err, shared.ErrInvalidPhone):
//		writeValidationError(w, map[string]string{"phone": "invalid_format"})
//	case errors.Is(err, authDomain.ErrPhoneAlreadyExists):
//		writeError(w, http.StatusConflict, "phone_already_exists", nil)
//	case errors.Is(err, billing.ErrPlanNotFound),
//		errors.Is(err, billing.ErrPlanInactive),
//		errors.Is(err, billing.ErrInvalidPlanCode):
//		writeError(w, http.StatusUnprocessableEntity, "invalid_plan", nil)
//	case errors.Is(err, authDomain.ErrOTPAlreadySent):
//		writeError(w, http.StatusTooManyRequests, "too_many_requests", nil)
//	default:
//		writeError(w, http.StatusInternalServerError, "internal_error", nil)
//	}
//}
