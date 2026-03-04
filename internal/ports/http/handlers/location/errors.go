package location

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
)

var locationErrors = httputil.ErrorMap{
	billing.ErrLimitExceeded:         {Code: http.StatusForbidden, Message: "limit_exceeded"},
	billing.ErrSubscriptionSuspended: {Code: http.StatusForbidden, Message: "subscription_suspended"},
	identity.ErrPermissionDenied:     {Code: http.StatusForbidden, Message: "permission_denied"},

	location.ErrLocationNotFound:     {Code: http.StatusNotFound, Message: "location_not_found"},
	location.ErrCategoryNotFound:     {Code: http.StatusBadRequest, Message: "category_not_found"},
	location.ErrNameRequired:         {Code: http.StatusBadRequest, Message: "name_required"},
	location.ErrInvalidScheduleType:  {Code: http.StatusBadRequest, Message: "invalid_schedule_type"},
	location.ErrScheduleTypeRequired: {Code: http.StatusBadRequest, Message: "schedule_type_required"},
	location.ErrCityRequired:         {Code: http.StatusBadRequest, Message: "city_required"},
	location.ErrInvalidLatitude:      {Code: http.StatusBadRequest, Message: "invalid_latitude"},
	location.ErrInvalidLongitude:     {Code: http.StatusBadRequest, Message: "invalid_longitude"},

	shared.ErrInvalidPhone: {Code: http.StatusBadRequest, Message: "invalid_phone"},
}
