package entity

import (
	"time"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/domain/enum"
)

type User struct {
	Phone     string    `json:"phone"`
	Role      enum.Role `json:"role"`
	PointCode string    `json:"point_code,omitempty"`
	Name      string    `json:"name,omitempty"`
	Surname   string    `json:"surname,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by,omitempty"`
}
