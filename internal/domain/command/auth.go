package command

import "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"

type RegisterStaffCommand struct {
	Phone     string
	PointCode string
	Role      enum.Role
}
type RegisterStaffConfirmCommand struct {
	Token    string
	Name     string
	Surname  string
	Password string
}
