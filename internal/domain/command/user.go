package command

import "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/users"

type UserUpdateCommand struct {
	Phone     string
	Name      *string
	Surname   *string
	UpdatedBy string
}

func (u *UserUpdateCommand) ToPatch() *users.UserUpdatePatch {
	patch := users.NewUserUpdatePatch().SetPhone(u.Phone)

	if u.Name != nil {
		patch.SetName(*u.Name)
	}
	if u.Surname != nil {
		patch.SetSurname(*u.Surname)
	}
	if u.UpdatedBy != "" {
		patch.SetUpdatedBy(u.UpdatedBy)
	}
	return patch.Build()
}
