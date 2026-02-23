package catalog

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type serviceModel struct {
	ID              uuid.UUID  `db:"id"`
	LocationID      uuid.UUID  `db:"location_id"`
	CategoryID      uuid.UUID  `db:"category_id"`
	Name            string     `db:"name"`
	Description     *string    `db:"description"`
	DurationMinutes int        `db:"duration_minutes"`
	Color           string     `db:"color"`
	SortOrder       int        `db:"sort_order"`
	Active          bool       `db:"active"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
}

func (m *serviceModel) toDomain() (*catalog.Service, error) {
	duration, err := shared.NewDuration(m.DurationMinutes)
	if err != nil {
		return nil, err
	}

	color, err := catalog.ParseColor(m.Color)
	if err != nil {
		return nil, err
	}

	return &catalog.Service{
		ID:              m.ID,
		LocationID:      m.LocationID,
		CategoryID:      m.CategoryID,
		Name:            m.Name,
		Description:     m.Description,
		DurationMinutes: duration,
		Color:           color,
		SortOrder:       m.SortOrder,
		Active:          m.Active,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		DeletedAt:       m.DeletedAt,
	}, nil
}

type employeeServiceModel struct {
	ID         uuid.UUID       `db:"id"`
	EmployeeID uuid.UUID       `db:"employee_id"`
	ServiceID  uuid.UUID       `db:"service_id"`
	Price      decimal.Decimal `db:"price"`
	Active     bool            `db:"active"`
	CreatedAt  time.Time       `db:"created_at"`
	UpdatedAt  *time.Time      `db:"updated_at"`
}

func (m *employeeServiceModel) toDomain() (*catalog.EmployeeService, error) {
	price, err := shared.NewMoney(m.Price)
	if err != nil {
		return nil, err
	}

	return &catalog.EmployeeService{
		ID:         m.ID,
		EmployeeID: m.EmployeeID,
		ServiceID:  m.ServiceID,
		Price:      price,
		Active:     m.Active,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}, nil
}
