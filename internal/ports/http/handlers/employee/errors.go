package employee

import (
	"net/http"

	authDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
)

var employeeErrors = httputil.ErrorMap{
	// action token / invite
	authDomain.ErrActionTokenNotFound: {Code: http.StatusNotFound, Message: "token_not_found"},
	authDomain.ErrInviteTokenExpired:  {Code: http.StatusGone, Message: "invite_token_expired"},
	authDomain.ErrInviteAlreadyExists: {Code: http.StatusTooManyRequests, Message: "invite_already_exists"},
	authDomain.ErrPhoneAlreadyExists:  {Code: http.StatusConflict, Message: "phone_already_exists"},

	// identity
	identity.ErrEmployeeNotFound:      {Code: http.StatusNotFound, Message: "employee_not_found"},
	identity.ErrEmployeeDeleted:       {Code: http.StatusGone, Message: "employee_deleted"},
	identity.ErrPermissionDenied:      {Code: http.StatusForbidden, Message: "permission_denied"},
	identity.ErrCannotModifySelf:      {Code: http.StatusForbidden, Message: "cannot_modify_self"},
	identity.ErrCannotSetOwnerRole:    {Code: http.StatusBadRequest, Message: "cannot_set_owner_role"},
	identity.ErrInvalidRole:           {Code: http.StatusBadRequest, Message: "invalid_role"},
	identity.ErrEmployeeNameRequired:  {Code: http.StatusBadRequest, Message: "name_required"},
	identity.ErrEmployeePhoneRequired: {Code: http.StatusBadRequest, Message: "phone_required"},

	// location
	location.ErrLocationNotFound: {Code: http.StatusNotFound, Message: "location_not_found"},
	location.ErrLocationInactive: {Code: http.StatusUnprocessableEntity, Message: "location_inactive"},
	location.ErrLocationDeleted:  {Code: http.StatusGone, Message: "location_deleted"},

	// media
	media.ErrAssetNotFound: {Code: http.StatusNotFound, Message: "avatar_not_found"},
	media.ErrAssetNotReady: {Code: http.StatusUnprocessableEntity, Message: "avatar_not_ready"},

	// billing
	billing.ErrSubscriptionSuspended: {Code: http.StatusForbidden, Message: "subscription_suspended"},

	// shared
	shared.ErrInvalidPhone: {Code: http.StatusBadRequest, Message: "invalid_phone"},
}
