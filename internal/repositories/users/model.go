package users

import (
	"encoding/json"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	"github.com/cockroachdb/errors"
)

type model struct {
	Phone         string     `db:"phone"`
	Password      string     `db:"password"`
	Role          string     `db:"role"`
	Name          string     `db:"name"`
	Surname       string     `db:"surname"`
	PointCode     *string    `db:"point_code"`
	NetworkCode   *string    `db:"network_code"`
	AvatarFileKey *string    `db:"avatar_file_key"` // Из JOIN с media таблицей
	Active        bool       `db:"active"`
	Schedule      []byte     `db:"schedule"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
	UpdatedBy     string     `db:"updated_by"`
}

func convert(e *user.User) (model, error) {
	m := model{
		Phone:       e.Phone,
		Password:    e.PasswordHash,
		Role:        e.Role.String(),
		Name:        e.Name,
		Surname:     e.Surname,
		PointCode:   e.PointCode,
		NetworkCode: e.NetworkCode,
		// AvatarFileKey НЕ сохраняем через users таблицу (обновляется через user_media)
		Active:    e.Active,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,
		UpdatedBy: e.UpdatedBy,
	}

	// Serialize Schedule to JSON via DTO
	// Важно: если Schedule пустой или nil, создаём пустой массив, а не null
	var schedulesDTO []staffScheduleDTO
	if e.Schedule != nil {
		schedulesDTO = schedulesToDTO(e.Schedule)
	} else {
		schedulesDTO = []staffScheduleDTO{}
	}

	bb, err := json.Marshal(schedulesDTO)
	if err != nil {
		return m, errors.Wrap(err, "failed to marshal schedule to json")
	}
	m.Schedule = bb

	return m, nil
}

func (m model) convert() (*user.User, error) {
	role, err := user.RoleString(m.Role)
	if err != nil {
		return nil, err
	}

	user := &user.User{
		Phone:        m.Phone,
		PasswordHash: m.Password,
		Role:         role,
		Name:         m.Name,
		Surname:      m.Surname,
		PointCode:    m.PointCode,
		NetworkCode:  m.NetworkCode,
		Avatar: &user.Avatar{
			FileKey: m.AvatarFileKey,
		}, // Маппим из JOIN
		// AvatarURL остается nil (заполняется в service при необходимости)
		Active:    m.Active,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
		UpdatedBy: m.UpdatedBy,
	}

	// Deserialize Schedule from JSON via DTO
	var scheduleDTOs schedulesDTO
	if err = json.Unmarshal(m.Schedule, &scheduleDTOs); err == nil {
		user.Schedule = scheduleDTOs.toEntity()
	}

	return user, nil
}

type models []model

func (m models) convert() ([]*user.User, error) {
	users := make([]*user.User, len(m))
	for i, model := range m {
		user, err := model.convert()
		if err != nil {
			return nil, err
		}
		users[i] = user
	}
	return users, nil
}
