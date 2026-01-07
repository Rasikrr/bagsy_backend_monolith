package users

import "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"

func convertToUpdatedUser(oldUser, newUser *entity.User) *entity.User {
	out := *oldUser

	if newUser.Name != "" {
		out.Name = newUser.Name
	}
	if newUser.Surname != "" {
		out.Surname = newUser.Surname
	}

	// остальные если надо добавим

	return &out
}
