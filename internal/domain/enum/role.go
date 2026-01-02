package enum

type Role int8

//go:generate enumer -type=Role -json -trimprefix Role -transform=snake -output role_enumer.go
const (
	RoleUser Role = iota
	RoleStaff
	RoleManager
	RoleNetManager
	RoleSelfOwner
	RoleAdmin
)

func (r Role) OneOf(roles ...Role) bool {
	for _, role := range roles {
		if role == r {
			return true
		}
	}
	return false
}
