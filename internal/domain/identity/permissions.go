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
