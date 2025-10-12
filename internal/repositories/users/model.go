// nolint: unused
package users

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
)

type model struct {
	Phone       string     `db:"phone"`
	Password    *string    `db:"password"`
	Role        string     `db:"role"`
	Name        *string    `db:"name"`
	Surname     *string    `db:"surname"`
	Active      bool       `db:"active"`
	PointCode   *string    `db:"point_code"`
	NetworkCode *string    `db:"network_code"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	UpdatedBy   *string    `db:"updated_by"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

type models []model

func (m model) convert() (*entity.User, error) {
	role, err := enum.RoleString(m.Role)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		Phone:       m.Phone,
		Role:        role,
		Name:        m.Name,
		Surname:     m.Surname,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		UpdatedBy:   m.UpdatedBy,
		PointCode:   m.PointCode,
		NetworkCode: m.NetworkCode,
		Active:      m.Active,
		Password:    m.Password,
		DeletedAt:   m.DeletedAt,
	}, nil
}

func (m models) convert() ([]*entity.User, error) {
	users := make([]*entity.User, len(m))
	for i, model := range m {
		user, err := model.convert()
		if err != nil {
			return nil, err
		}
		users[i] = user
	}
	return users, nil
}

func convert(user *entity.User) (*model, error) {
	role, err := enum.RoleString(user.Role.String())
	if err != nil {
		return nil, err
	}

	out := &model{
		Phone:       user.Phone,
		Role:        role.String(),
		Password:    user.Password,
		Name:        user.Name,
		Surname:     user.Surname,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		UpdatedBy:   user.UpdatedBy,
		PointCode:   user.PointCode,
		NetworkCode: user.NetworkCode,
	}
	now := time.Now().UTC()
	if out.CreatedAt.IsZero() {
		out.CreatedAt = now
	}
	if out.UpdatedAt != nil && out.UpdatedAt.IsZero() {
		out.UpdatedAt = &now
	}
	return out, nil
}
