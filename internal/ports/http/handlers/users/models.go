package users

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/hash"
)

//go:generate easyjson -all models.go

type userResponse struct {
	Name      string     `json:"name"`
	Surname   string     `json:"surname"`
	Role      string     `json:"role"`
	PointCode string     `json:"point_code"`
	Phone     string     `json:"phone"`
	Active    bool       `json:"active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy *string    `json:"updated_by"`
}

func convertUserResponse(u *entity.User) userResponse {
	return userResponse{
		Name:      u.Name,
		Role:      u.Role.String(),
		Surname:   u.Surname,
		Phone:     u.Phone,
		Active:    u.Active,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		UpdatedBy: u.UpdatedBy,
	}
}

type updateRequest struct {
	Name     *string `json:"name"`
	Surname  *string `json:"surname"`
	Password *string `json:"password"`
}

func (u *updateRequest) toParams() (users.UpdateParams, error) {
	params := users.UpdateParams{
		Name:    u.Name,
		Surname: u.Surname,
	}
	if u.Password != nil {
		hashedPassword, err := hash.Password(*u.Password)
		if err != nil {
			return users.UpdateParams{}, err
		}
		params.Password = &hashedPassword
	}
	return params, nil
}
