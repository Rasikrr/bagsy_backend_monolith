package policy

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/booking"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
)

type Policy struct{}

func New() *Policy {
	return &Policy{}
}

func (p *Policy) CanCancelAppointment(orgCtx *access.OrgContext, appt *booking.Appointment) error {
	if !orgCtx.Subscription.Status.CanOperate() {
		return billing.ErrSubscriptionSuspended
	}

	// 1. Проверка организации (базовая безопасность)
	if !appt.BelongsTo(orgCtx.Organization.ID) {
		return identity.ErrPermissionDenied
	}

	// 2. Проверка ролей
	switch {
	case orgCtx.Employee.Role.IsOwner():
		// Владелец может отменять всё в рамках организации
		return nil

	case orgCtx.Employee.Role.IsManager():
		// Менеджер может отменять только в рамках своей локации
		if !appt.BelongsToLocation(orgCtx.Employee.LocationID) {
			return identity.ErrPermissionDenied
		}
		return nil

	case orgCtx.Employee.Role.IsStaff():
		// Мастер может отменять только свои записи
		if !appt.BelongsToEmployee(orgCtx.Employee.ID) {
			return identity.ErrPermissionDenied
		}
		return nil

	default:
		return identity.ErrPermissionDenied
	}
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

func (p *Policy) CanListEmployees(orgCtx *access.OrgContext, filter *identity.EmployeeFilter) error {
	if !orgCtx.Subscription.Status.CanOperate() {
		return billing.ErrSubscriptionSuspended
	}

	switch {
	case orgCtx.Employee.Role.IsOwner():
		// Owner может видеть всех сотрудников организации
		return nil

	case orgCtx.Employee.Role.IsManager():
		// Manager видит только сотрудников своей локации
		locID := orgCtx.Employee.LocationID
		filter.LocationID = &locID
		return nil

	default:
		return identity.ErrPermissionDenied
	}
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
