package identity

import (
	"time"

	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────────
// ChangeType — причина создания записи в истории.
// ─────────────────────────────────────────────────────────────────

type ChangeType string

const (
	ChangeTypeHired     ChangeType = "hired"
	ChangeTypePromotion ChangeType = "promotion"
	ChangeTypeDemotion  ChangeType = "demotion"
	ChangeTypeTransfer  ChangeType = "transfer"
)

func (c ChangeType) IsValid() bool {
	switch c {
	case ChangeTypeHired, ChangeTypePromotion,
		ChangeTypeDemotion, ChangeTypeTransfer:
		return true
	}
	return false
}

func (c ChangeType) String() string {
	return string(c)
}

// ─────────────────────────────────────────────────────────────────
// WorkHistory — запись в истории работы сотрудника.
//
// Каждая строка — период работы в конкретной роли на конкретной точке.
// При изменении (повышение, перевод, увольнение) текущая запись
// закрывается (EndedAt), создаётся новая.
// ─────────────────────────────────────────────────────────────────

type WorkHistory struct {
	ID             uuid.UUID
	EmployeeID     uuid.UUID
	OrganizationID uuid.UUID
	LocationID     uuid.UUID
	Role           Role
	StartedAt      time.Time
	EndedAt        *time.Time
	ChangeType     ChangeType
	Comment        *string
	CreatedAt      time.Time
}

func NewWorkHistory(
	employeeID uuid.UUID,
	organizationID uuid.UUID,
	locationID uuid.UUID,
	role Role,
	changeType ChangeType,
	comment *string,
) *WorkHistory {
	now := time.Now()
	return &WorkHistory{
		ID:             uuid.New(),
		EmployeeID:     employeeID,
		OrganizationID: organizationID,
		LocationID:     locationID,
		Role:           role,
		StartedAt:      now,
		ChangeType:     changeType,
		Comment:        comment,
		CreatedAt:      now,
	}
}

// End закрывает текущий период.
func (w *WorkHistory) End(at time.Time) {
	w.EndedAt = &at
}

// IsActive возвращает true если период ещё не закрыт.
func (w *WorkHistory) IsActive() bool {
	return w.EndedAt == nil
}
