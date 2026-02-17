package identity

import (
	"strings"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────────
// Aggregate Root: Employee
// ─────────────────────────────────────────────────────────────────

type Employee struct {
	ID           uuid.UUID
	Phone        shared.Phone
	PasswordHash string
	FirstName    string
	LastName     *string

	AvatarID *uuid.UUID

	OrganizationID uuid.UUID
	LocationID     uuid.UUID

	Role        Role
	Permissions Permissions

	Active    bool
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

type CreateEmployeeParams struct {
	Phone          shared.Phone
	FirstName      string
	LastName       *string
	OrganizationID uuid.UUID
	LocationID     uuid.UUID
	Role           Role
	Permissions    Permissions
}

type CreateOwnerParams struct {
	Phone          shared.Phone
	FirstName      string
	LastName       *string
	OrganizationID uuid.UUID
	Permissions    Permissions
}

func NewOwnerEmployee(params CreateOwnerParams) (*Employee, error) {
	if err := validateEmployeeName(params.FirstName); err != nil {
		return nil, err
	}
	if params.Phone.IsEmpty() {
		return nil, ErrEmployeePhoneRequired
	}

	return &Employee{
		ID:             uuid.New(),
		Phone:          params.Phone,
		FirstName:      strings.TrimSpace(params.FirstName),
		LastName:       params.LastName,
		OrganizationID: params.OrganizationID,
		Role:           RoleOwner,
		Permissions:    params.Permissions,
		Active:         true,
		CreatedAt:      time.Now(),
	}, nil
}

func NewEmployee(params CreateEmployeeParams) (*Employee, error) {
	if err := validateEmployeeName(params.FirstName); err != nil {
		return nil, err
	}
	if params.Phone.IsEmpty() {
		return nil, ErrEmployeePhoneRequired
	}
	if !params.Role.IsValid() {
		return nil, ErrInvalidRole
	}

	return &Employee{
		ID:             uuid.New(),
		Phone:          params.Phone,
		FirstName:      strings.TrimSpace(params.FirstName),
		LastName:       params.LastName,
		OrganizationID: params.OrganizationID,
		LocationID:     params.LocationID,
		Role:           params.Role,
		Permissions:    params.Permissions,
		Active:         true,
		CreatedAt:      time.Now(),
	}, nil
}

// ─────────────────────────────────────────────────────────────────
// Business Methods
// ─────────────────────────────────────────────────────────────────

func (e *Employee) UpdateProfile(firstName string, lastName *string) error {
	if e.IsDeleted() {
		return ErrEmployeeDeleted
	}
	if err := validateEmployeeName(firstName); err != nil {
		return err
	}

	e.FirstName = strings.TrimSpace(firstName)
	e.LastName = lastName
	e.touch()
	return nil
}

func (e *Employee) ChangeAvatar(avatarID uuid.UUID) error {
	if e.IsDeleted() {
		return ErrEmployeeDeleted
	}
	if e.AvatarID != nil && *e.AvatarID == avatarID {
		return nil
	}
	e.AvatarID = &avatarID
	e.touch()
	return nil
}

func (e *Employee) ChangeRole(newRole Role) error {
	if e.IsDeleted() {
		return ErrEmployeeDeleted
	}
	if !newRole.IsValid() {
		return ErrInvalidRole
	}

	e.Role = newRole
	e.touch()
	return nil
}

func (e *Employee) SetPermissions(permissions Permissions) error {
	if e.IsDeleted() {
		return ErrEmployeeDeleted
	}

	e.Permissions = permissions
	e.touch()
	return nil
}

func (e *Employee) Transfer(locationID uuid.UUID) error {
	if e.IsDeleted() {
		return ErrEmployeeDeleted
	}

	e.LocationID = locationID
	e.touch()
	return nil
}

func (e *Employee) Activate() error {
	if e.IsDeleted() {
		return ErrEmployeeDeleted
	}
	e.Active = true
	e.touch()
	return nil
}

func (e *Employee) Deactivate() error {
	if e.IsDeleted() {
		return ErrEmployeeDeleted
	}
	e.Active = false
	e.touch()
	return nil
}

func (e *Employee) Delete() error {
	if e.IsDeleted() {
		return nil
	}
	now := time.Now()
	e.DeletedAt = &now
	e.Active = false
	return nil
}

func (e *Employee) SetPassword(hash string) {
	e.PasswordHash = hash
	e.touch()
}

// ─────────────────────────────────────────────────────────────────
// Query Methods
// ─────────────────────────────────────────────────────────────────

func (e *Employee) IsDeleted() bool {
	return e.DeletedAt != nil
}

func (e *Employee) IsActive() bool {
	return e.Active && !e.IsDeleted()
}

func (e *Employee) FullName() string {
	if e.LastName == nil {
		return e.FirstName
	}
	return e.FirstName + " " + *e.LastName
}

// ─────────────────────────────────────────────────────────────────
// Private Methods
// ─────────────────────────────────────────────────────────────────

func (e *Employee) touch() {
	now := time.Now()
	e.UpdatedAt = &now
}

func validateEmployeeName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrEmployeeNameRequired
	}
	return nil
}
