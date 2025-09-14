package enum

type Role int8

//go:generate enumer -type=Role -json -trimprefix Role -transform=snake -output role_enumer.go
const (
	RoleAdmin Role = iota
	RoleNetManager
	RoleManager
	RoleStaff
	RoleUser
	RoleModerator
	RoleSelfOwner
)
