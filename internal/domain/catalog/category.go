package catalog

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────────
// Aggregate Root: ServiceCategory
// ─────────────────────────────────────────────────────────────────

type ServiceCategory struct {
	ID                 uuid.UUID
	LocationCategoryID uuid.UUID
	ParentID           *uuid.UUID
	Name               string
	SortOrder          int
	CreatedAt          time.Time
	UpdatedAt          *time.Time
}

func NewServiceCategory(
	locationCategoryID uuid.UUID,
	parentID *uuid.UUID,
	name string,
) (*ServiceCategory, error) {
	if err := validateServiceCategoryName(name); err != nil {
		return nil, err
	}
	return &ServiceCategory{
		ID:                 uuid.New(),
		LocationCategoryID: locationCategoryID,
		ParentID:           parentID,
		Name:               name,
		CreatedAt:          time.Now(),
	}, nil
}

func (c *ServiceCategory) UpdateName(name string) error {
	if err := validateServiceCategoryName(name); err != nil {
		return err
	}

	c.Name = name
	c.touch()

	return nil
}

func (c *ServiceCategory) ChangeParent(parentID *uuid.UUID) error {
	// Нельзя сделать родителем самого себя
	if parentID != nil && *parentID == c.ID {
		return ErrCategorySelfParent
	}

	c.ParentID = parentID
	c.touch()

	return nil
}

func (c *ServiceCategory) IsRoot() bool {
	return c.ParentID == nil
}

func (c *ServiceCategory) touch() {
	now := time.Now()
	c.UpdatedAt = &now
}

func validateServiceCategoryName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrCategoryNameRequired
	}
	return nil
}
