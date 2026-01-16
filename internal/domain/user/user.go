package user

import (
	"time"
)

type User struct {
	Phone        string
	PasswordHash string
	Role         Role
	Name         string
	Surname      string
	PointCode    *string
	NetworkCode  *string
	Avatar       *Avatar
	Active       bool
	Schedule     Schedule
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time
	UpdatedBy    string
}

// DetachFromLocation - отвязка юзера от места работы
func (u *User) DetachFromLocation() {
	u.Schedule = nil
	u.PointCode = nil
	u.NetworkCode = nil
}

func (u *User) IsAssignedToLocation() bool {
	return u.PointCode != nil || u.NetworkCode != nil
}

type Schedule []*ScheduleElement

type ScheduleElement struct {
	WeekDay int
	Open    time.Time
	Close   time.Time
	AllDay  bool
	Comment string
}

type Avatar struct {
	FileKey *string
	URL     string
}

func (a *Avatar) GetURL() string {
	return a.URL
}
