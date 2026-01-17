package forms

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/form"
)

type model struct {
	ID          int       `db:"id"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	Role        string    `db:"role"`
	Phone       string    `db:"phone"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

func toModel(form *form.Form) *model {
	return &model{
		ID:          form.ID,
		FirstName:   form.FirstName,
		LastName:    form.LastName,
		Role:        form.Role,
		Phone:       form.Phone,
		Description: form.Description,
		CreatedAt:   form.CreatedAt,
	}
}
