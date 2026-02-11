package catalog

import (
	"strings"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────────
// Aggregate Root: Service
// ─────────────────────────────────────────────────────────────────

type Service struct {
	ID              uuid.UUID
	LocationID      uuid.UUID
	CategoryID      uuid.UUID
	Name            string
	Description     *string
	DurationMinutes shared.Duration
	SortOrder       int
	Active          bool
	CreatedAt       time.Time
	UpdatedAt       *time.Time
	DeletedAt       *time.Time
}

type CreateServiceParams struct {
	LocationID      uuid.UUID
	CategoryID      uuid.UUID
	Name            string
	Description     *string
	DurationMinutes shared.Duration
}

func NewService(params CreateServiceParams) (*Service, error) {
	if err := validateServiceName(params.Name); err != nil {
		return nil, err
	}

	return &Service{
		ID:              uuid.New(),
		LocationID:      params.LocationID,
		CategoryID:      params.CategoryID,
		Name:            strings.TrimSpace(params.Name),
		Description:     params.Description,
		DurationMinutes: params.DurationMinutes,
		Active:          true,
		CreatedAt:       time.Now(),
	}, nil
}

// ─────────────────────────────────────────────────────────────────
// Business Methods
// ─────────────────────────────────────────────────────────────────

func (s *Service) UpdateInfo(
	name string,
	description *string,
	duration shared.Duration,
) error {
	if s.IsDeleted() {
		return ErrServiceDeleted
	}

	if err := validateServiceName(name); err != nil {
		return err
	}

	s.Name = strings.TrimSpace(name)
	s.Description = description
	s.DurationMinutes = duration
	s.touch()

	return nil
}

func (s *Service) ChangeCategory(categoryID uuid.UUID) error {
	if s.IsDeleted() {
		return ErrServiceDeleted
	}

	s.CategoryID = categoryID
	s.touch()

	return nil
}

func (s *Service) ChangeSortOrder(order int) {
	s.SortOrder = order
	s.touch()
}

func (s *Service) Activate() error {
	if s.IsDeleted() {
		return ErrServiceDeleted
	}

	s.Active = true
	s.touch()

	return nil
}

func (s *Service) Deactivate() error {
	if s.IsDeleted() {
		return ErrServiceDeleted
	}

	s.Active = false
	s.touch()

	return nil
}

func (s *Service) Delete() error {
	if s.IsDeleted() {
		return nil
	}

	now := time.Now()
	s.DeletedAt = &now
	s.Active = false

	return nil
}

// ─────────────────────────────────────────────────────────────────
// Query Methods
// ─────────────────────────────────────────────────────────────────

func (s *Service) IsDeleted() bool {
	return s.DeletedAt != nil
}

func (s *Service) IsActive() bool {
	return s.Active && !s.IsDeleted()
}

func (s *Service) BelongsToLocation(locationID uuid.UUID) bool {
	return s.LocationID == locationID
}

// ─────────────────────────────────────────────────────────────────
// Private
// ─────────────────────────────────────────────────────────────────

func (s *Service) touch() {
	now := time.Now()
	s.UpdatedAt = &now
}

func validateServiceName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrServiceNameRequired
	}
	return nil
}
