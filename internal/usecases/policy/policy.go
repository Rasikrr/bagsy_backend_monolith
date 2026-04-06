package policy

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/booking"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/google/uuid"
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
	if orgCtx.Employee.ID == targetEmp.ID && !orgCtx.Employee.Role.IsOwner() {
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

// CanUnassignEmployee проверяет право на отвязку сотрудника от точки. Только owner.
func (p *Policy) CanUnassignEmployee(orgCtx *access.OrgContext, targetEmp *identity.Employee) error {
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

// CanViewLocation проверяет право на просмотр одной локации.
// Owner — любую, Manager/Staff — только свою.
func (p *Policy) CanViewLocation(orgCtx *access.OrgContext, locationID uuid.UUID) error {
	if !orgCtx.Subscription.Status.CanOperate() {
		return billing.ErrSubscriptionSuspended
	}

	if orgCtx.Employee.Role.IsOwner() {
		return nil
	}

	if orgCtx.Employee.LocationID != locationID {
		return identity.ErrPermissionDenied
	}
	return nil
}

// CanViewLocations проверяет право на просмотр списка локаций. Только owner.
func (p *Policy) CanViewLocations(orgCtx *access.OrgContext) error {
	if !orgCtx.Subscription.Status.CanOperate() {
		return billing.ErrSubscriptionSuspended
	}
	if !orgCtx.Employee.Role.IsOwner() {
		return identity.ErrPermissionDenied
	}
	return nil
}

// CanCreateService проверяет право на создание услуги в локации.
// Owner — в любой локации org, Manager — только в своей.
func (p *Policy) CanCreateService(orgCtx *access.OrgContext, locationID uuid.UUID) error {
	if !orgCtx.Subscription.Status.CanOperate() {
		return billing.ErrSubscriptionSuspended
	}

	switch {
	case orgCtx.Employee.Role.IsOwner():
		return nil
	case orgCtx.Employee.Role.IsManager():
		if orgCtx.Employee.LocationID != locationID {
			return identity.ErrPermissionDenied
		}
		return nil
	default:
		return identity.ErrPermissionDenied
	}
}

// CanManageService проверяет право на обновление/удаление услуги в локации.
// Owner — в любой локации org, Manager — только в своей.
func (p *Policy) CanManageService(orgCtx *access.OrgContext, locationID uuid.UUID) error {
	if !orgCtx.Subscription.Status.CanOperate() {
		return billing.ErrSubscriptionSuspended
	}

	switch {
	case orgCtx.Employee.Role.IsOwner():
		return nil
	case orgCtx.Employee.Role.IsManager():
		if orgCtx.Employee.LocationID != locationID {
			return identity.ErrPermissionDenied
		}
		return nil
	default:
		return identity.ErrPermissionDenied
	}
}

// CanCreateEmployeeService проверяет право на привязку услуги к сотруднику.
// Owner — любого, Manager — staff в своей локации.
func (p *Policy) CanCreateEmployeeService(orgCtx *access.OrgContext, targetEmp *identity.Employee) error {
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

// CanManageLocationSchedule — Owner: any location in org. Manager: own location + CanManageLocationSchedule permission.
func (p *Policy) CanManageLocationSchedule(orgCtx *access.OrgContext, locationID uuid.UUID) error {
	if !orgCtx.Subscription.Status.CanOperate() {
		return billing.ErrSubscriptionSuspended
	}

	switch {
	case orgCtx.Employee.Role.IsOwner():
		return nil
	case orgCtx.Employee.Role.IsManager():
		if !orgCtx.Employee.Permissions.CanManageLocationSchedule {
			return identity.ErrPermissionDenied
		}
		if orgCtx.Employee.LocationID != locationID {
			return identity.ErrPermissionDenied
		}
		return nil
	default:
		return identity.ErrPermissionDenied
	}
}

// CanViewLocationSchedule — Owner: any. Others: own location.
func (p *Policy) CanViewLocationSchedule(orgCtx *access.OrgContext, locationID uuid.UUID) error {
	if !orgCtx.Subscription.Status.CanOperate() {
		return billing.ErrSubscriptionSuspended
	}

	if orgCtx.Employee.Role.IsOwner() {
		return nil
	}

	if orgCtx.Employee.LocationID != locationID {
		return identity.ErrPermissionDenied
	}
	return nil
}

// CanManageEmployeeSchedule — Owner: any employee in org. Manager: staff in own location. Staff: self only.
func (p *Policy) CanManageEmployeeSchedule(orgCtx *access.OrgContext, targetEmp *identity.Employee) error {
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
		if !targetEmp.Role.IsStaff() && !orgCtx.Employee.Permissions.CanProvideServices {
			return identity.ErrPermissionDenied
		}
		if targetEmp.LocationID == nil || *targetEmp.LocationID != orgCtx.Employee.LocationID {
			return identity.ErrPermissionDenied
		}
		return nil
	case orgCtx.Employee.Role.IsStaff():
		if orgCtx.Employee.ID != targetEmp.ID {
			return identity.ErrPermissionDenied
		}
		return nil
	default:
		return identity.ErrPermissionDenied
	}
}

// CanViewEmployeeSchedule — Owner: any. Manager: own location employees. Staff: self only.
func (p *Policy) CanViewEmployeeSchedule(orgCtx *access.OrgContext, targetEmp *identity.Employee) error {
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
		if targetEmp.LocationID == nil || *targetEmp.LocationID != orgCtx.Employee.LocationID {
			return identity.ErrPermissionDenied
		}
		return nil
	case orgCtx.Employee.Role.IsStaff():
		if orgCtx.Employee.ID != targetEmp.ID {
			return identity.ErrPermissionDenied
		}
		return nil
	default:
		return identity.ErrPermissionDenied
	}
}

// CanManageLocation проверяет право на обновление/удаление локации. Только owner.
func (p *Policy) CanManageLocation(orgCtx *access.OrgContext, loc *location.Location) error {
	if !orgCtx.Subscription.Status.CanOperate() {
		return billing.ErrSubscriptionSuspended
	}
	if !loc.BelongsTo(orgCtx.Organization.ID) {
		return identity.ErrPermissionDenied
	}
	if !orgCtx.Employee.Role.IsOwner() {
		return identity.ErrPermissionDenied
	}
	return nil
}

func (p *Policy) CanCreateDirectBooking(orgCtx *access.OrgContext, locationID uuid.UUID, targetEmployeeID uuid.UUID) error {
	if !orgCtx.Subscription.Status.CanOperate() {
		return billing.ErrSubscriptionSuspended
	}
	switch {
	case orgCtx.Employee.Role.IsOwner():
		return nil
	case orgCtx.Employee.Role.IsManager():
		if orgCtx.Employee.LocationID != locationID {
			return identity.ErrPermissionDenied
		}
		return nil
	case orgCtx.Employee.Role.IsStaff():
		if orgCtx.Employee.ID != targetEmployeeID {
			return identity.ErrPermissionDenied
		}
		return nil
	default:
		return identity.ErrPermissionDenied
	}
}

// CanUpdateOrganization проверяет право на обновление профиля организации. Только owner.
func (p *Policy) CanUpdateOrganization(orgCtx *access.OrgContext) error {
	if !orgCtx.Subscription.Status.CanOperate() {
		return billing.ErrSubscriptionSuspended
	}
	if !orgCtx.Employee.Role.IsOwner() {
		return identity.ErrPermissionDenied
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
