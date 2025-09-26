// nolint: unused
package users

import (
	"time"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/domain/enum"
)

type model struct {
	Phone     string     `db:"phone"`
	Role      string     `db:"role"`
	Name      string     `db:"name"`
	Surname   string     `db:"surname"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	UpdatedBy *string    `db:"updated_by"`
	PointCode *string    `db:"point_code"`
	Active    bool       `db:"active"`
	Password  *string    `db:"password"`
}

type models []model

func (m model) convert() (*entity.User, error) {
	role, err := enum.RoleString(m.Role)
	if err != nil {
		return nil, err
	}

	var pointCode string
	if m.PointCode != nil {
		pointCode = *m.PointCode
	}

	return &entity.User{
		Phone:     m.Phone,
		Role:      role,
		Name:      m.Name,
		Surname:   m.Surname,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		UpdatedBy: m.UpdatedBy,
		PointCode: pointCode,
		Active:    m.Active,
		Password:  m.Password,
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

	return &model{
		Phone:     user.Phone,
		Role:      role.String(),
		Name:      user.Name,
		Surname:   user.Surname,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		UpdatedBy: user.UpdatedBy,
		PointCode: &user.PointCode,
	}, nil
}
