package policy

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
)

type Policy struct{}

func New() *Policy {
	return &Policy{}
}

func (p *Policy) CanCreateLocation(orgCtx *access.OrgContext, currentCount int) error {
	if !orgCtx.Subscription.Status.CanOperate() {
		return billing.ErrSubscriptionSuspended
	}

	if !orgCtx.Employee.Role.IsOwner() {
		return identity.ErrPermissionDenied
	}

	if !orgCtx.Plan.Capabilities.CanUse(billing.ResourceMaxLocations, currentCount) {
		return billing.ErrLimitExceeded
	}

	return nil
}
