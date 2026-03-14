package schedule

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	domainSchedule "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
)

var scheduleErrors = httputil.ErrorMap{
	billing.ErrSubscriptionSuspended: {Code: http.StatusForbidden, Message: "subscription_suspended"},
	identity.ErrPermissionDenied:     {Code: http.StatusForbidden, Message: "permission_denied"},
	identity.ErrEmployeeNotFound:     {Code: http.StatusNotFound, Message: "employee_not_found"},
	location.ErrLocationNotFound:     {Code: http.StatusNotFound, Message: "location_not_found"},
	location.ErrLocationInactive:     {Code: http.StatusUnprocessableEntity, Message: "location_inactive"},

	domainSchedule.ErrInvalidSlotType:  {Code: http.StatusBadRequest, Message: "invalid_slot_type"},
	domainSchedule.ErrInvalidTimeRange: {Code: http.StatusBadRequest, Message: "invalid_time_range"},
}
