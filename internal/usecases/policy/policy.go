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

// CanManageEmployee проверяет право на activate, deactivate.
func (p *Policy) CanManageEmployee(orgCtx *access.OrgContext, targetEmp *identity.Employee) error {
	if !orgCtx.Subscription.Status.CanOperate() {
		return billing.ErrSubscriptionSuspended
	}
	if orgCtx.Employee.ID == targetEmp.ID {
		return identity.ErrCannotModifySelf
	}
	if orgCtx.Organization.ID != targetEmp.OrganizationID {
		return identity.ErrPermissionDenied
	}

	switch {
	case orgCtx.Employee.Role.IsOwner():
		return nil
	case orgCtx.Employee.Role.IsManager():
		if !targetEmp.Role.IsStaff() {
			return identity.ErrPermissionDenied
		}
		if (targetEmp.LocationID == nil) ||
			(*targetEmp.LocationID != orgCtx.Employee.LocationID) {
			return identity.ErrPermissionDenied
		}
		return nil
	default:
		return identity.ErrPermissionDenied
	}
}

// CanTransferEmployee проверяет право на перевод сотрудника в другую локацию. Только owner.
func (p *Policy) CanTransferEmployee(orgCtx *access.OrgContext, targetEmp *identity.Employee) error {
	if !orgCtx.Subscription.Status.CanOperate() {
		return billing.ErrSubscriptionSuspended
	}
	if orgCtx.Employee.ID == targetEmp.ID {
		return identity.ErrCannotModifySelf
	}
	if orgCtx.Organization.ID != targetEmp.OrganizationID {
		return identity.ErrPermissionDenied
	}
	if !orgCtx.Employee.Role.IsOwner() {
		return identity.ErrPermissionDenied
	}
	return nil
}

// CanChangeRole проверяет право на смену роли. Только owner, нельзя менять себя.
func (p *Policy) CanChangeRole(orgCtx *access.OrgContext, targetEmp *identity.Employee) error {
	if !orgCtx.Subscription.Status.CanOperate() {
		return billing.ErrSubscriptionSuspended
	}
	if orgCtx.Employee.ID == targetEmp.ID {
		return identity.ErrCannotModifySelf
	}
	if orgCtx.Organization.ID != targetEmp.OrganizationID {
		return identity.ErrPermissionDenied
	}
	if !orgCtx.Employee.Role.IsOwner() {
		return identity.ErrPermissionDenied
	}
	return nil
}

// CanChangePermissions проверяет право на смену permissions.
// Owner — любого, включая себя (намеренно: может убрать себя из бронирования).
// Manager — staff в своей локации.
func (p *Policy) CanChangePermissions(orgCtx *access.OrgContext, targetEmp *identity.Employee) error {
	if !orgCtx.Subscription.Status.CanOperate() {
		return billing.ErrSubscriptionSuspended
	}
	if orgCtx.Organization.ID != targetEmp.OrganizationID {
		return identity.ErrPermissionDenied
	}

	switch {
	case orgCtx.Employee.Role.IsOwner():
		return nil
	case orgCtx.Employee.Role.IsManager():
		if !targetEmp.Role.IsStaff() {
			return identity.ErrPermissionDenied
		}
		if targetEmp.LocationID == nil || *targetEmp.LocationID != orgCtx.Employee.LocationID {
			return identity.ErrPermissionDenied
		}
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
