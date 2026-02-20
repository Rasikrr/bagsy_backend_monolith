package employee

import (
	"net/http"

	authDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
)

var employeeErrors = util.ErrorMap{
	// action token / invite
	authDomain.ErrActionTokenNotFound: {Code: http.StatusNotFound, Message: "token_not_found"},
	authDomain.ErrInviteTokenExpired:  {Code: http.StatusGone, Message: "invite_token_expired"},
	authDomain.ErrInviteAlreadyExists: {Code: http.StatusTooManyRequests, Message: "invite_already_exists"},
	authDomain.ErrPhoneAlreadyExists:  {Code: http.StatusConflict, Message: "phone_already_exists"},

	// identity
	identity.ErrPermissionDenied:      {Code: http.StatusForbidden, Message: "permission_denied"},
	identity.ErrInvalidRole:           {Code: http.StatusBadRequest, Message: "invalid_role"},
	identity.ErrEmployeeNameRequired:  {Code: http.StatusBadRequest, Message: "name_required"},
	identity.ErrEmployeePhoneRequired: {Code: http.StatusBadRequest, Message: "phone_required"},

	// billing
	billing.ErrSubscriptionSuspended: {Code: http.StatusForbidden, Message: "subscription_suspended"},

	// shared
	shared.ErrInvalidPhone: {Code: http.StatusBadRequest, Message: "invalid_phone"},
}
