package billing

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
)

var billingErrors = httputil.ErrorMap{
	billing.ErrSubscriptionNotFound:         {Code: http.StatusNotFound, Message: "subscription_not_found"},
	billing.ErrSubscriptionSuspended:        {Code: http.StatusForbidden, Message: "subscription_suspended"},
	billing.ErrSubscriptionActive:           {Code: http.StatusConflict, Message: "subscription_already_active"},
	billing.ErrInvalidStatusTransition:      {Code: http.StatusUnprocessableEntity, Message: "invalid_status_transition"},
	billing.ErrInvalidBillingCycle:          {Code: http.StatusBadRequest, Message: "invalid_billing_cycle"},
	billing.ErrNotActiveForCancellation:     {Code: http.StatusUnprocessableEntity, Message: "not_active_for_cancellation"},
	billing.ErrCancellationAlreadyRequested: {Code: http.StatusConflict, Message: "cancellation_already_requested"},
	billing.ErrNoCancellationToUndo:         {Code: http.StatusUnprocessableEntity, Message: "no_cancellation_to_undo"},
	billing.ErrPlanNotFound:                 {Code: http.StatusNotFound, Message: "plan_not_found"},
	identity.ErrPermissionDenied:            {Code: http.StatusForbidden, Message: "permission_denied"},
}
