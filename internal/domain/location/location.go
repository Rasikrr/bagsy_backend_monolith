package location

import (
	"strings"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────────
// Aggregate Root: Location
// ─────────────────────────────────────────────────────────────────

type Location struct {
	ID                  uuid.UUID
	OrganizationID      uuid.UUID
	CategoryID          uuid.UUID
	Name                string
	Description         *string
	Phone               *shared.Phone
	Slug                shared.Slug
	Address             *Address
	Coordinates         *Coordinates
	Active              bool
	ScheduleType        ScheduleType
	SlotDurationMinutes shared.Duration
	CreatedAt           time.Time
	UpdatedAt           *time.Time
	DeletedAt           *time.Time
}

type CreateLocationParams struct {
	OrganizationID      uuid.UUID
	CategoryID          uuid.UUID
	Name                string
	Description         *string
	Phone               *shared.Phone
	Address             *Address
	Coordinates         *Coordinates
	ScheduleType        ScheduleType
	SlotDurationMinutes shared.Duration
}

func NewLocation(params CreateLocationParams) (*Location, error) {
	if err := validateLocationName(params.Name); err != nil {
		return nil, err
	}

	if !params.ScheduleType.IsValid() {
		return nil, ErrInvalidScheduleType
	}

	slug, err := shared.NewSlug(params.Name)
	if err != nil {
		return nil, err
	}

	return &Location{
		ID:                  uuid.New(),
		OrganizationID:      params.OrganizationID,
		CategoryID:          params.CategoryID,
		Name:                params.Name,
		Description:         params.Description,
		Phone:               params.Phone,
		Slug:                slug,
		Address:             params.Address,
		Coordinates:         params.Coordinates,
		ScheduleType:        params.ScheduleType,
		SlotDurationMinutes: params.SlotDurationMinutes,
		Active:              true,
		CreatedAt:           time.Now(),
	}, nil
}

// ─────────────────────────────────────────────────────────────────
// Business Methods
// ─────────────────────────────────────────────────────────────────

func (l *Location) UpdateInfo(name string, description *string) error {
	if l.IsDeleted() {
		return ErrLocationDeleted
	}

	if err := validateLocationName(name); err != nil {
		return err
	}

	l.Name = name
	l.Description = description
	l.touch()

	return nil
}

func (l *Location) ChangeCategory(categoryID uuid.UUID) error {
	if l.IsDeleted() {
		return ErrLocationDeleted
	}

	l.CategoryID = categoryID
	l.touch()

	return nil
}

func (l *Location) ChangePhone(phone *shared.Phone) error {
	if l.IsDeleted() {
		return ErrLocationDeleted
	}

	l.Phone = phone
	l.touch()

	return nil
}

func (l *Location) SetAddress(address *Address, coordinates *Coordinates) error {
	if l.IsDeleted() {
		return ErrLocationDeleted
	}

	l.Address = address
	l.Coordinates = coordinates
	l.touch()

	return nil
}

func (l *Location) ChangeSlug(new shared.Slug) error {
	if l.IsDeleted() {
		return ErrLocationDeleted
	}
	if l.Slug.IsEqual(new) {
		return nil
	}
	l.Slug = new
	l.touch()
	return nil
}

func (l *Location) ChangeScheduleType(scheduleType ScheduleType) error {
	if l.IsDeleted() {
		return ErrLocationDeleted
	}

	if !scheduleType.IsValid() {
		return ErrInvalidScheduleType
	}

	l.ScheduleType = scheduleType
	l.touch()

	return nil
}

func (l *Location) ChangeSlotDuration(minutesDur shared.Duration) error {
	if l.IsDeleted() {
		return ErrLocationDeleted
	}

	l.SlotDurationMinutes = minutesDur
	l.touch()

	return nil
}

func (l *Location) Activate() error {
	if l.IsDeleted() {
		return ErrLocationDeleted
	}

	l.Active = true
	l.touch()

	return nil
}

func (l *Location) Deactivate() error {
	if l.IsDeleted() {
		return ErrLocationDeleted
	}

	l.Active = false
	l.touch()

	return nil
}

func (l *Location) Delete() error {
	if l.IsDeleted() {
		return nil
	}

	now := time.Now()
	l.DeletedAt = &now
	l.Active = false

	return nil
}

// ─────────────────────────────────────────────────────────────────
// Query Methods
// ─────────────────────────────────────────────────────────────────

func (l *Location) IsDeleted() bool {
	return l.DeletedAt != nil
}

func (l *Location) CanOperate() bool {
	return l.Active && !l.IsDeleted()
}

func (l *Location) BelongsTo(organizationID uuid.UUID) bool {
	return l.OrganizationID == organizationID
}

func (l *Location) HasFixedSchedule() bool {
	return l.ScheduleType == ScheduleTypeFixed
}

func (l *Location) HasMixedSchedule() bool {
	return l.ScheduleType == ScheduleTypeMixed
}

// ─────────────────────────────────────────────────────────────────
// Private Methods
// ─────────────────────────────────────────────────────────────────

func (l *Location) touch() {
	now := time.Now()
	l.UpdatedAt = &now
}

func validateLocationName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrNameRequired
	}
	return nil
}
