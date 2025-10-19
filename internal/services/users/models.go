package users

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/users"
)

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
