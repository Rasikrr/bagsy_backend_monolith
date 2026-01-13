package entity

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
)

type User struct {
	Phone         string
	Password      string
	Role          enum.Role
	Name          string
	Surname       string
	PointCode     *string
	NetworkCode   *string
	AvatarFileKey *string // file_key из media таблицы (через JOIN user_media → media), персистится
	Active        bool
	Schedule      []StaffSchedule
	CreatedAt     time.Time
	UpdatedAt     *time.Time
	DeletedAt     *time.Time
	UpdatedBy     string
}

type StaffSchedule struct {
	WeekDay int
	Open    time.Time
	Close   time.Time
	AllDay  bool
	Comment string
}
