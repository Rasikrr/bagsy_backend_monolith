package auth

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
)

type RegisterStaffCommand struct {
	Phone       string
	Name        string
	Surname     string
	PointCode   string
	NetworkCode string
	Role        user.Role
}
type RegisterStaffConfirmCommand struct {
	Token    string
	Password string
}

type ChangePasswordConfirmCommand RegisterStaffConfirmCommand

type RegisterManagementCommand struct {
	Name                string
	Surname             string
	Phone               string
	Password            string
	Role                user.Role
	NetworkRegisterInfo *RegisterNetworkInfo
}

type RegisterNetworkInfo struct {
	Name        string
	Description string
}
