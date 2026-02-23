package auth

import (
	"net/http"

	authDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
)

var authErrors = httputil.ErrorMap{
	// auth domain
	authDomain.ErrOTPExpired:           {Code: http.StatusGone, Message: "otp_expired"},
	authDomain.ErrOTPInvalid:           {Code: http.StatusUnprocessableEntity, Message: "otp_invalid"},
	authDomain.ErrOTPAlreadySent:       {Code: http.StatusTooManyRequests, Message: "otp_already_sent"},
	authDomain.ErrTooManyAttempts:      {Code: http.StatusTooManyRequests, Message: "too_many_attempts"},
	authDomain.ErrPhoneAlreadyExists:   {Code: http.StatusConflict, Message: "phone_already_exists"},
	authDomain.ErrRegistrationExpired:  {Code: http.StatusGone, Message: "registration_expired"},
	authDomain.ErrRefreshTokenNotFound: {Code: http.StatusUnauthorized, Message: "token_not_found"},
	authDomain.ErrRefreshTokenExpired:  {Code: http.StatusUnauthorized, Message: "token_expired"},
	authDomain.ErrActionTokenNotFound:  {Code: http.StatusNotFound, Message: "token_not_found"},
	authDomain.ErrEmployeeInactive:     {Code: http.StatusForbidden, Message: "employee_inactive"},

	// identity domain — в контексте auth "not found" означает неверные креды
	identity.ErrEmployeeNotFound: {Code: http.StatusUnauthorized, Message: "invalid_credentials"},

	// billing domain
	billing.ErrPlanNotFound:    {Code: http.StatusBadRequest, Message: "plan_not_found"},
	billing.ErrInvalidPlanCode: {Code: http.StatusBadRequest, Message: "invalid_plan_code"},

	// shared domain
	shared.ErrInvalidPhone: {Code: http.StatusBadRequest, Message: "invalid_phone"},
}
