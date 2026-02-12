package identity

type Role string

const (
	RoleOwner   Role = "owner"
	RoleManager Role = "manager"
	RoleStaff   Role = "staff"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleOwner, RoleManager, RoleStaff:
		return true
	}
	return false
}

func (r Role) String() string {
	return string(r)
}

func ParseRole(s string) (Role, error) {
	role := Role(s)
	if !role.IsValid() {
		return "", ErrInvalidRole
	}
	return role, nil
}

func (r Role) IsOwner() bool {
	return r == RoleOwner
}

func (r Role) IsManager() bool {
	return r == RoleManager
}

func (r Role) IsStaff() bool {
	return r == RoleStaff
}
