package users

import (
	"errors"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/users"
	coreErr "github.com/Rasikrr/core/errors"
	"github.com/samber/lo"
)

type GetParams struct {
	NetworkCode *string
	PointCode   *string
	Phones      []string
	Roles       []enum.Role
}

func (g GetParams) validate(by *entity.Session) error {
	if by.GetRole() == enum.RoleAdmin {
		return nil
	}
	if g.NetworkCode != nil {
		if by.GetRole().OneOf(enum.RoleManager, enum.RoleNetManager, enum.RoleSelfOwner) {
			if by.GetNetworkCode() != *g.NetworkCode {
				return coreErr.ErrForbidden.Wrap(errors.New("invalid network code"))
			}
		}
	}
	if g.PointCode != nil {
		if by.GetRole().OneOf(enum.RoleManager, enum.RoleSelfOwner) {
			if by.GetPointCode() != *g.PointCode {
				return coreErr.ErrForbidden.Wrap(errors.New("invalid point code"))
			}
		}
	}
	return nil
}

func (g *GetParams) convert() users.GetParams {
	return users.GetParams{
		NetworkCode: g.NetworkCode,
		PointCode:   g.PointCode,
		Phones:      g.Phones,
		Roles: lo.Map(g.Roles, func(r enum.Role, _ int) string {
			return r.String()
		}),
	}
}

type UpdateParams struct {
	Name     *string
	Surname  *string
	Password *string
}

func (u *UpdateParams) ToPatch(phone string) *users.UserUpdatePatch {
	patch := users.NewUserUpdatePatch().SetPhones(phone)

	if u.Name != nil {
		patch.SetName(*u.Name)
	}
	if u.Surname != nil {
		patch.SetSurname(*u.Surname)
	}
	if u.Password != nil {
		patch.SetPassword(*u.Password)
	}
	return patch.Build()
}
