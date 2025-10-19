package users

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/hash"
	"github.com/samber/lo"
)

//go:generate easyjson -all models.go

type userResponse struct {
	Name        *string    `json:"name,omitempty"`
	Surname     *string    `json:"surname,omitempty"`
	Role        string     `json:"role"`
	PointCode   *string    `json:"point_code,omitempty"`
	NetworkCode *string    `json:"network_code,omitempty"`
	Phone       string     `json:"phone"`
	Active      bool       `json:"active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	UpdatedBy   *string    `json:"updated_by,omitempty"`
}

type userListResponse struct {
	Users []userResponse `json:"users"`
}

func convertUserListResponse(users []*entity.User) userListResponse {
	return userListResponse{
		Users: lo.Map(users, func(u *entity.User, _ int) userResponse {
			return convertUserResponse(u)
		}),
	}
}

func convertUserResponse(u *entity.User) userResponse {
	return userResponse{
		Name:        u.Name,
		Role:        u.Role.String(),
		Surname:     u.Surname,
		PointCode:   u.PointCode,
		NetworkCode: u.NetworkCode,
		Phone:       u.Phone,
		Active:      u.Active,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		UpdatedBy:   u.UpdatedBy,
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
