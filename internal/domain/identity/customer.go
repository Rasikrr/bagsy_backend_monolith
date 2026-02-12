package identity

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────────
// Aggregate Root: Customer
// ─────────────────────────────────────────────────────────────────

type Customer struct {
	ID        uuid.UUID
	Phone     Phone
	FirstName string
	LastName  *string
	BirthDate *time.Time

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func NewCustomer(phone Phone, firstName string, lastName *string) (*Customer, error) {
	if phone.IsEmpty() {
		return nil, ErrCustomerPhoneRequired
	}
	if strings.TrimSpace(firstName) == "" {
		return nil, ErrCustomerNameRequired
	}

	return &Customer{
		ID:        uuid.New(),
		Phone:     phone,
		FirstName: strings.TrimSpace(firstName),
		LastName:  lastName,
		CreatedAt: time.Now(),
	}, nil
}

// ─────────────────────────────────────────────────────────────────
// Business Methods
// ─────────────────────────────────────────────────────────────────

func (c *Customer) UpdateProfile(firstName string, lastName *string, birthDate *time.Time) error {
	if c.IsDeleted() {
		return ErrCustomerDeleted
	}
	if strings.TrimSpace(firstName) == "" {
		return ErrCustomerNameRequired
	}

	c.FirstName = strings.TrimSpace(firstName)
	c.LastName = lastName
	c.BirthDate = birthDate
	c.touch()
	return nil
}

func (c *Customer) Delete() error {
	if c.IsDeleted() {
		return nil
	}
	now := time.Now()
	c.DeletedAt = &now
	return nil
}

// ─────────────────────────────────────────────────────────────────
// Query Methods
// ─────────────────────────────────────────────────────────────────

func (c *Customer) IsDeleted() bool {
	return c.DeletedAt != nil
}

func (c *Customer) FullName() string {
	if c.LastName == nil {
		return c.FirstName
	}
	return c.FirstName + " " + *c.LastName
}

// ─────────────────────────────────────────────────────────────────
// Private Methods
// ─────────────────────────────────────────────────────────────────

func (c *Customer) touch() {
	now := time.Now()
	c.UpdatedAt = &now
}
