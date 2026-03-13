package catalog

import (
	"fmt"
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

func fromServiceDomain(s *catalog.Service) serviceModel {
	return serviceModel{
		ID:              s.ID,
		LocationID:      s.LocationID,
		CategoryID:      s.CategoryID,
		Name:            s.Name,
		Description:     s.Description,
		DurationMinutes: s.DurationMinutes.Minutes(),
		Color:           string(s.Color),
		SortOrder:       s.SortOrder,
		Active:          s.Active,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
		DeletedAt:       s.DeletedAt,
	}
}

type serviceWithPricesModel struct {
	serviceModel
	MinPrice *decimal.Decimal `db:"min_price"`
	MaxPrice *decimal.Decimal `db:"max_price"`
}

func (m *serviceWithPricesModel) toDomain() (*catalog.Service, error) {
	svc, err := m.serviceModel.toDomain()
	if err != nil {
		return nil, err
	}

	var (
		money shared.Money
	)
	if m.MinPrice != nil {
		money, err = shared.NewMoney(*m.MinPrice)
		if err != nil {
			return nil, fmt.Errorf("parse min_price: %w", err)
		}
		svc.MinPrice = &money
	}
	if m.MaxPrice != nil {
		money, err = shared.NewMoney(*m.MaxPrice)
		if err != nil {
			return nil, fmt.Errorf("parse max_price: %w", err)
		}
		svc.MaxPrice = &money
	}
	return svc, nil
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

func fromEmployeeServiceDomain(es *catalog.EmployeeService) employeeServiceModel {
	return employeeServiceModel{
		ID:         es.ID,
		EmployeeID: es.EmployeeID,
		ServiceID:  es.ServiceID,
		Price:      es.Price.Amount(),
		Active:     es.Active,
		CreatedAt:  es.CreatedAt,
		UpdatedAt:  es.UpdatedAt,
	}
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

type serviceCategoryModel struct {
	ID                 uuid.UUID  `db:"id"`
	LocationCategoryID uuid.UUID  `db:"location_category_id"`
	ParentID           *uuid.UUID `db:"parent_id"`
	Name               string     `db:"name"`
	SortOrder          int        `db:"sort_order"`
	CreatedAt          time.Time  `db:"created_at"`
}

func (m *serviceCategoryModel) toDomain() *catalog.ServiceCategory {
	return &catalog.ServiceCategory{
		ID:                 m.ID,
		LocationCategoryID: m.LocationCategoryID,
		ParentID:           m.ParentID,
		Name:               m.Name,
		SortOrder:          m.SortOrder,
		CreatedAt:          m.CreatedAt,
	}
}
