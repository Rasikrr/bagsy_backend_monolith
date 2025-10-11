package entity

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
)

type User struct {
	Phone     string     `json:"phone"`
	Password  *string    `json:"password,omitempty"`
	Role      enum.Role  `json:"role"`
	Name      *string    `json:"name,omitempty"`
	Surname   *string    `json:"surname,omitempty"`
	Active    bool       `json:"active"`
	PointCode *string    `json:"point_code,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	UpdatedBy *string    `json:"updated_by"`
}

func NewCustomerUser(phone string) *User {
	return &User{
		Phone:     phone,
		Role:      enum.RoleUser,
		CreatedAt: time.Now(),
		Active:    true,
	}
}
