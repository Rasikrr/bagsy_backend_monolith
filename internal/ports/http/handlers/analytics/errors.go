package analytics

import (
	"net/http"

	domainAnalytics "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/analytics"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
)

var analyticsErrors = httputil.ErrorMap{
	domainAnalytics.ErrAccessDenied:  {Code: http.StatusForbidden, Message: "access_denied"},
	domainAnalytics.ErrNotFound:      {Code: http.StatusNotFound, Message: "not_found"},
	domainAnalytics.ErrInvalidPeriod: {Code: http.StatusUnprocessableEntity, Message: "invalid_period"},

	billing.ErrSubscriptionSuspended: {Code: http.StatusForbidden, Message: "subscription_suspended"},
}
