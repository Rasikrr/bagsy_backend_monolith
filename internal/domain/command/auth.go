package command

import "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"

type RegisterStaffCommand struct {
	Phone       string
	PointCode   string
	NetworkCode string
	Role        enum.Role
}
type RegisterStaffConfirmCommand struct {
	Token    string
	Name     string
	Surname  string
	Password string
}

type RegisterManagementCommand struct {
	Name                string
	Surname             string
	Phone               string
	Password            string
	Role                enum.Role
	NetworkRegisterInfo *NetworkRegisterInfo

	AuthCode string
	Attempts int
}

type NetworkRegisterInfo struct {
	Name        string
	Description string
}
