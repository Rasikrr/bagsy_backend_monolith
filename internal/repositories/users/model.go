package users

import (
	"encoding/json"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/cockroachdb/errors"
)

type model struct {
	Phone       string     `db:"phone"`
	Password    string     `db:"password"`
	Role        string     `db:"role"`
	Name        string     `db:"name"`
	Surname     string     `db:"surname"`
	PointCode   *string    `db:"point_code"`
	NetworkCode *string    `db:"network_code"`
	Active      bool       `db:"active"`
	Schedule    []byte     `db:"schedule"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	UpdatedBy   string     `db:"updated_by"`
}

func convert(e *entity.User) (model, error) {
	m := model{
		Phone:       e.Phone,
		Password:    e.Password,
		Role:        e.Role.String(),
		Name:        e.Name,
		Surname:     e.Surname,
		PointCode:   e.PointCode,
		NetworkCode: e.NetworkCode,
		Active:      e.Active,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		DeletedAt:   e.DeletedAt,
		UpdatedBy:   e.UpdatedBy,
	}

	// Serialize Schedule to JSON if not nil
	if e.Schedule != nil {
		scheduleBytes, err := json.Marshal(e.Schedule)
		if err != nil {
			return model{}, errors.Wrap(err, "failed to marshal schedule to json")
		}
		m.Schedule = scheduleBytes
	}

	return m, nil
}

func (m model) convert() (*entity.User, error) {
	role, err := enum.RoleString(m.Role)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Phone:       m.Phone,
		Password:    m.Password,
		Role:        role,
		Name:        m.Name,
		Surname:     m.Surname,
		PointCode:   m.PointCode,
		NetworkCode: m.NetworkCode,
		Active:      m.Active,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
		UpdatedBy:   m.UpdatedBy,
	}

	// Deserialize Schedule from JSON if not empty
	if len(m.Schedule) > 0 {
		var schedule entity.StaffSchedule
		if err = json.Unmarshal(m.Schedule, &schedule); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal schedule from json")
		}
		user.Schedule = &schedule
	}

	return user, nil
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
