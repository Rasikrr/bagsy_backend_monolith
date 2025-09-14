package entity

import "time"

type Bagsy struct {
	ID        int       `json:"id"`
	Time      time.Time `json:"time"`
	PointCode string    `json:"point_code"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by,omitempty"`
}
