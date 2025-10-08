package users

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
)

//go:generate easyjson -all models.go

type updateRequest struct {
	Name     *string `json:"name"`
	Surname  *string `json:"surname"`
	Password *string `json:"password"`
}

func (u *updateRequest) toParams() users.UpdateParams {
	return users.UpdateParams{
		Name:     u.Name,
		Surname:  u.Surname,
		Password: u.Password,
	}
}
