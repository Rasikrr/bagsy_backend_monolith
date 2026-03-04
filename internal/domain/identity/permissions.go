package identity

type Permissions struct {
	CanProvideServices        bool
	CanManageLocationSchedule bool
}

func NewPermissions(canProvideServices, canManageSchedule bool) Permissions {
	return Permissions{
		CanProvideServices:        canProvideServices,
		CanManageLocationSchedule: canManageSchedule,
	}
}

func DefaultPermissions() Permissions {
	return Permissions{
		CanProvideServices:        false,
		CanManageLocationSchedule: false,
	}
}

// DefaultPermissionsForRole возвращает дефолтные permissions для роли.
func DefaultPermissionsForRole(role Role) Permissions {
	switch role {
	case RoleOwner:
		return NewPermissions(true, true)
	case RoleManager:
		return NewPermissions(false, true)
	case RoleStaff:
		return NewPermissions(true, false)
	default:
		return DefaultPermissions()
	}
}
