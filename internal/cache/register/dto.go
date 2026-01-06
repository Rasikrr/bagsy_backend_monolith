package register

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/command"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
)

type registerManagementDTO struct {
	Name                string               `json:"name"`
	Surname             string               `json:"surname"`
	Phone               string               `json:"phone"`
	Password            string               `json:"password"`
	Role                string               `json:"role"`
	NetworkRegisterInfo *networkRegisterInfo `json:"network_register_info"`

	AuthCode string `json:"auth_code"`
	Attempts int    `json:"attempts"`
}

type networkRegisterInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *registerManagementDTO) toDomain() *command.RegisterManagementCommand {
	role, _ := enum.RoleString(r.Role)
	return &command.RegisterManagementCommand{
		Name:     r.Name,
		Surname:  r.Surname,
		Phone:    r.Phone,
		Password: r.Password,
		Role:     role,
		NetworkRegisterInfo: &command.NetworkRegisterInfo{
			Name:        r.NetworkRegisterInfo.Name,
			Description: r.NetworkRegisterInfo.Description,
		},
		AuthCode: r.AuthCode,
		Attempts: r.Attempts,
	}
}

func toRegisterManagementDTO(req *command.RegisterManagementCommand) *registerManagementDTO {
	return &registerManagementDTO{
		Name:     req.Name,
		Surname:  req.Surname,
		Phone:    req.Phone,
		Password: req.Password,
		Role:     req.Role.String(),
		NetworkRegisterInfo: &networkRegisterInfo{
			Name:        req.NetworkRegisterInfo.Name,
			Description: req.NetworkRegisterInfo.Description,
		},
		AuthCode: req.AuthCode,
		Attempts: req.Attempts,
	}
}

type registerStaffDTO struct {
	Phone       string `json:"phone"`
	PointCode   string `json:"point_code"`
	NetworkCode string `json:"network_code"`
	Role        string `json:"role"`
}

func toRegisterStaffDTO(c *command.RegisterStaffCommand) *registerStaffDTO {
	return &registerStaffDTO{
		Phone:       c.Phone,
		PointCode:   c.PointCode,
		NetworkCode: c.NetworkCode,
		Role:        c.Role.String(),
	}
}

func (r *registerStaffDTO) toDomain() *command.RegisterStaffCommand {
	role, _ := enum.RoleString(r.Role)
	return &command.RegisterStaffCommand{
		Phone:       r.Phone,
		PointCode:   r.PointCode,
		NetworkCode: r.NetworkCode,
		Role:        role,
	}
}
