package register

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	authS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"
)

type registerManagementStateDTO struct {
	Command  *registerManagementCommandDTO `json:"command"`
	AuthCode string                        `json:"auth_code"`
	Attempts int                           `json:"attempts"`
}

type registerManagementCommandDTO struct {
	Name                string               `json:"name"`
	Surname             string               `json:"surname"`
	Phone               string               `json:"phone"`
	Password            string               `json:"password"`
	Role                string               `json:"role"`
	NetworkRegisterInfo *networkRegisterInfo `json:"network_register_info"`
}

type networkRegisterInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *registerManagementStateDTO) toDomain() *authS.ManagementRegistrationState {
	role, _ := user.RoleString(r.Command.Role)
	return &authS.ManagementRegistrationState{
		Command: &auth.RegisterManagementCommand{
			Name:     r.Command.Name,
			Surname:  r.Command.Surname,
			Phone:    r.Command.Phone,
			Password: r.Command.Password,
			Role:     role,
			NetworkRegisterInfo: &auth.RegisterNetworkInfo{
				Name:        r.Command.NetworkRegisterInfo.Name,
				Description: r.Command.NetworkRegisterInfo.Description,
			},
		},
		AuthCode: r.AuthCode,
		Attempts: r.Attempts,
	}
}

func toRegisterManagementDTO(req *authS.ManagementRegistrationState) *registerManagementStateDTO {
	return &registerManagementStateDTO{
		Command: &registerManagementCommandDTO{
			Name:     req.Command.Name,
			Surname:  req.Command.Surname,
			Phone:    req.Command.Phone,
			Password: req.Command.Password,
			Role:     req.Command.Role.String(),
			NetworkRegisterInfo: &networkRegisterInfo{
				Name:        req.Command.NetworkRegisterInfo.Name,
				Description: req.Command.NetworkRegisterInfo.Description,
			},
		},
		AuthCode: req.AuthCode,
		Attempts: req.Attempts,
	}
}

type registerStaffDTO struct {
	Phone       string `json:"phone"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	PointCode   string `json:"point_code"`
	NetworkCode string `json:"network_code"`
	Role        string `json:"role"`
}

func toRegisterStaffDTO(c *auth.RegisterStaffCommand) *registerStaffDTO {
	return &registerStaffDTO{
		Phone:       c.Phone,
		Name:        c.Name,
		Surname:     c.Surname,
		PointCode:   c.PointCode,
		NetworkCode: c.NetworkCode,
		Role:        c.Role.String(),
	}
}

func (r *registerStaffDTO) toDomain() *auth.RegisterStaffCommand {
	role, _ := user.RoleString(r.Role)
	return &auth.RegisterStaffCommand{
		Name:        r.Name,
		Surname:     r.Surname,
		Phone:       r.Phone,
		PointCode:   r.PointCode,
		NetworkCode: r.NetworkCode,
		Role:        role,
	}
}
