package user

import (
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/google/uuid"
)

type UpdateUserCommand struct {
	Name     string
	Surname  string
	AvatarID *uuid.UUID
}

type CreateOwnerCommand struct {
	Name        string
	Surname     string
	Password    string
	Phone       string
	Role        Role
	NetworkCode string
}

type PromoteToOwnerCommand struct {
	Name        string
	Surname     string
	Password    string
	Role        Role
	NetworkCode string
}

func (p *PromoteToOwnerCommand) ToPromoteNewLocation() *PromoteToNewLocationCommand {
	return &PromoteToNewLocationCommand{
		Name:        p.Name,
		Surname:     p.Surname,
		Password:    p.Password,
		Role:        p.Role,
		NetworkCode: p.NetworkCode,
	}
}

type PromoteToStaffCommand struct {
	Name        string
	Surname     string
	Password    string
	Role        Role
	NetworkCode string
	PointCode   string
}

func (p *PromoteToStaffCommand) ToPromoteNewLocation() *PromoteToNewLocationCommand {
	return &PromoteToNewLocationCommand{
		Name:        p.Name,
		Surname:     p.Surname,
		Password:    p.Password,
		Role:        p.Role,
		NetworkCode: p.NetworkCode,
		PointCode:   &p.PointCode,
	}
}

type PromoteToNewLocationCommand struct {
	Name        string
	Surname     string
	Password    string
	Role        Role
	NetworkCode string
	PointCode   *string
}

func (p *PromoteToNewLocationCommand) Validate() error {
	if p.NetworkCode == "" {
		return domainErr.NewInvalidInputError(
			"staff and owners could not have empty network code",
			nil,
		)
	}

	switch p.Role {
	case RoleNetManager, RoleSelfOwner:

	case RoleStaff:
		if p.PointCode == nil || *p.PointCode == "" {
			return domainErr.NewInvalidInputError(
				"point code required for staff role",
				nil,
			)
		}

	default:
		return domainErr.NewInvalidInputError(
			"unsupported role for promotion",
			nil,
		).WithDetail("role", p.Role.String())
	}

	return nil
}

type CreateStaffCommand struct {
	Name        string
	Surname     string
	Password    string
	Phone       string
	Role        Role
	PointCode   string
	NetworkCode string
}

type CreateUserCommand struct {
	Name    string
	Surname string
	Phone   string
}
