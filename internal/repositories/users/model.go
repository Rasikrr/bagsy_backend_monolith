package users

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
)

type model struct {
	Phone       string     `db:"phone"`
	Password    string     `db:"password"`
	Role        string     `db:"role"`
	Name        string     `db:"name"`
	Surname     string     `db:"surname"`
	PointCode   *string    `db:"point_code"`
	NetworkCode *string    `db:"network_code"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	UpdatedBy   string     `db:"updated_by"`
}

func convert(e *entity.User) model {
	return model{
		Phone:       e.Phone,
		Password:    e.Password,
		Role:        e.Role.String(),
		Name:        e.Name,
		Surname:     e.Surname,
		PointCode:   e.PointCode,
		NetworkCode: e.NetworkCode,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		DeletedAt:   e.DeletedAt,
		UpdatedBy:   e.UpdatedBy,
	}
}

func (m model) convert() (*entity.User, error) {
	role, err := enum.RoleString(m.Role)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		Phone:       m.Phone,
		Password:    m.Password,
		Role:        role,
		Name:        m.Name,
		Surname:     m.Surname,
		PointCode:   m.PointCode,
		NetworkCode: m.NetworkCode,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
		UpdatedBy:   m.UpdatedBy,
	}, nil
}

type models []model

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
