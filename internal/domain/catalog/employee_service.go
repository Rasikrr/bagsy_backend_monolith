package catalog

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────────
// Aggregate Root: EmployeeService
// ─────────────────────────────────────────────────────────────────

type EmployeeService struct {
	ID         uuid.UUID
	EmployeeID uuid.UUID
	ServiceID  uuid.UUID
	Price      shared.Money
	Active     bool
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

func NewEmployeeService(
	employeeID uuid.UUID,
	serviceID uuid.UUID,
	price shared.Money,
) (*EmployeeService, error) {
	return &EmployeeService{
		ID:         uuid.New(),
		EmployeeID: employeeID,
		ServiceID:  serviceID,
		Price:      price,
		Active:     true,
		CreatedAt:  time.Now(),
	}, nil
}

func (es *EmployeeService) UpdatePrice(price shared.Money) error {
	es.Price = price
	es.touch()
	return nil
}

func (es *EmployeeService) Activate() {
	es.Active = true
	es.touch()
}

func (es *EmployeeService) Deactivate() {
	es.Active = false
	es.touch()
}

func (es *EmployeeService) IsActive() bool {
	return es.Active
}

func (es *EmployeeService) touch() {
	now := time.Now()
	es.UpdatedAt = &now
}
