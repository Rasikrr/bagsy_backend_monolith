package location

import (
	"strings"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID
	Slug      shared.Slug
	Name      string
	SortOrder int
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func NewCategory(name string) (*Category, error) {
	if err := validateCategoryName(name); err != nil {
		return nil, err
	}

	// 1. Генерируем слаг из имени прямо здесь
	slug, err := shared.NewSlug(name)
	if err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)

	return &Category{
		ID:        uuid.New(),
		Slug:      slug,
		Name:      name,
		CreatedAt: time.Now(),
	}, nil
}

// ─────────────────────────────────────────────────────────────────
// Business Methods
// ─────────────────────────────────────────────────────────────────

// UpdateName меняет ТОЛЬКО отображаемое имя.
// Слаг остается старым, чтобы не ломать SEO-ссылки.
func (c *Category) UpdateName(name string) error {
	if err := validateCategoryName(name); err != nil {
		return err
	}

	// Если имя не изменилось - ничего не делаем (и не трогаем UpdatedAt)
	if c.Name == strings.TrimSpace(name) {
		return nil
	}

	c.Name = strings.TrimSpace(name)
	c.touch()
	return nil
}

// ChangeSlug используется, если мы РЕАЛЬНО хотим сменить URL.
// Это должно быть отдельным действием в админке с предупреждением "Ссылки могут сломаться".
func (c *Category) ChangeSlug(newSlug shared.Slug) {
	if c.Slug == newSlug {
		return
	}
	c.Slug = newSlug
	c.touch()
}

// Reorder меняет порядок сортировки.
func (c *Category) Reorder(sortOrder int) {
	if c.SortOrder == sortOrder {
		return
	}
	c.SortOrder = sortOrder
	c.touch()
}

// ─────────────────────────────────────────────────────────────────
// Private Methods
// ─────────────────────────────────────────────────────────────────

func (c *Category) touch() {
	now := time.Now()
	c.UpdatedAt = &now
}

func validateCategoryName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrCategoryNameRequired
	}
	return nil
}
