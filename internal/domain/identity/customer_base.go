package identity

import (
	"time"

	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────────
// Aggregate Root: CustomerBase
// ─────────────────────────────────────────────────────────────────

type CustomerBase struct {
	ID             uuid.UUID
	CustomerID     uuid.UUID
	OrganizationID uuid.UUID

	FirstName string
	LastName  *string
	BirthDate *time.Time
	Gender    *Gender

	Notes []CustomerNote

	CreatedAt time.Time
	UpdatedAt *time.Time
}

// CustomerNote is an Entity within CustomerBase aggregate
type CustomerNote struct {
	ID             uuid.UUID
	CustomerBaseID uuid.UUID
	AuthorID       uuid.UUID
	Note           string
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}

type CreateCustomerBaseParams struct {
	CustomerID     uuid.UUID
	OrganizationID uuid.UUID
	FirstName      string
	LastName       *string
	BirthDate      *time.Time
	Gender         *Gender
}

func NewCustomerBase(params CreateCustomerBaseParams) *CustomerBase {
	return &CustomerBase{
		ID:             uuid.New(),
		CustomerID:     params.CustomerID,
		OrganizationID: params.OrganizationID,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
		BirthDate:      params.BirthDate,
		Gender:         params.Gender,
		Notes:          make([]CustomerNote, 0),
		CreatedAt:      time.Now(),
	}
}

// ─────────────────────────────────────────────────────────────────
// Business Methods
// ─────────────────────────────────────────────────────────────────

func (cb *CustomerBase) UpdateInfo(
	firstName string,
	lastName *string,
	birthDate *time.Time,
	gender *Gender,
) {
	cb.FirstName = firstName
	cb.LastName = lastName
	cb.BirthDate = birthDate
	cb.Gender = gender
	cb.touch()
}

func (cb *CustomerBase) AddNote(authorID uuid.UUID, content string) {
	note := CustomerNote{
		ID:             uuid.New(),
		CustomerBaseID: cb.ID,
		AuthorID:       authorID,
		Note:           content,
		CreatedAt:      time.Now(),
	}
	cb.Notes = append(cb.Notes, note)
	cb.touch()
}

// ─────────────────────────────────────────────────────────────────
// Private Methods
// ─────────────────────────────────────────────────────────────────

func (cb *CustomerBase) touch() {
	now := time.Now()
	cb.UpdatedAt = &now
}
