package entity

import "time"

type Network struct {
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	UpdatedBy   *string    `json:"updated_by,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
