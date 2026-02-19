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

func (p *Policy) CanInviteEmployee(orgCtx *access.OrgContext, targetRole identity.Role, currentCount int) error {
	if !orgCtx.Subscription.Status.CanOperate() {
		return billing.ErrSubscriptionSuspended
	}

	switch {
	case orgCtx.Employee.Role.IsOwner():
		if targetRole.IsOwner() {
			return identity.ErrPermissionDenied
		}
	case orgCtx.Employee.Role.IsManager():
		if !targetRole.IsStaff() {
			return identity.ErrPermissionDenied
		}
	default:
		return identity.ErrPermissionDenied
	}

	if !orgCtx.Plan.Capabilities.CanUse(billing.ResourceMaxEmployees, currentCount) {
		return billing.ErrLimitExceeded
	}

	return nil
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
