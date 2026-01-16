package form

import "time"

type Form struct {
	ID          int
	FirstName   string
	LastName    string
	Role        string
	Phone       string
	Description string
	CreatedAt   time.Time
}
