package entity

import "time"

type Network struct {
	Code        string
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
	UpdatedBy   string
}
