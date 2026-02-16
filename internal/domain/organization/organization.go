package organization

import (
	"strings"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────────
// Aggregate Root: Organization
// ─────────────────────────────────────────────────────────────────

type Organization struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Slug        shared.Slug
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}

// NewStubOrganization создает техническую организацию-контейнер.
// Используется при регистрации первой точки.
func NewStubOrganization() (*Organization, error) {
	org := &Organization{
		ID:        uuid.New(),
		Active:    true,
		CreatedAt: time.Now(),
	}
	return org, nil
}

// SetupProfile заполняет данные организации.
// Обязательно вызывать перед добавлением второй локации.
func (o *Organization) SetupProfile(name string, description *string) error {
	if o.IsDeleted() {
		return ErrOrganizationDeleted
	}
	if err := validateOrganizationName(name); err != nil {
		return err
	}
	slug, err := shared.NewSlug(name)
	if err != nil {
		return err
	}
	o.Name = name
	o.Description = description
	o.Slug = slug

	o.touch()

	return nil
}

// UpdateInfo - Безопасное обновление информации
// Используется в 99% случаев. Ссылки не ломаются.
func (o *Organization) UpdateInfo(name string, description *string) error {
	if o.IsDeleted() {
		return ErrOrganizationDeleted
	}
	if err := validateOrganizationName(name); err != nil {
		return err
	}
	name = strings.TrimSpace(name)

	// Мы меняем отображаемое имя, но URL (Code) оставляем старым!
	o.Name = name
	o.Description = description
	o.touch()

	return nil
}

func (o *Organization) IsProfileComplete() bool {
	return o.Name != "" && !o.Slug.IsEmpty()
}

// ChangeSlug - Опасное обновление слага (Ребрендинг)
// Вызывать только если пользователь нажал кнопку "Изменить ссылку" и подтвердил,
// что понимает риски (старые ссылки перестанут работать).
func (o *Organization) ChangeSlug(newSlug shared.Slug) error {
	if o.IsDeleted() {
		return ErrOrganizationDeleted
	}

	// Если организация еще "черновик" (без слага, имени), то менять нечего.
	// Нужно сначала вызвать SetupProfile.
	if !o.IsProfileComplete() {
		return ErrOrganizationProfileIncomplete
	}

	// Если слаг такой же - игнорируем
	if o.Slug.IsEqual(newSlug) {
		return nil
	}

	o.Slug = newSlug
	o.touch()

	return nil
}

func (o *Organization) IsActive() bool {
	return o.Active
}

func (o *Organization) Activate() error {
	if o.IsDeleted() {
		return ErrOrganizationDeleted
	}

	if o.Active {
		return nil
	}

	o.Active = true
	o.touch()

	return nil
}

func (o *Organization) Deactivate() error {
	if o.IsDeleted() {
		return ErrOrganizationDeleted
	}

	if !o.Active {
		return nil
	}

	o.Active = false
	o.touch()

	return nil
}

func (o *Organization) CanOperate() bool {
	return o.Active && !o.IsDeleted()
}

func (o *Organization) Delete() error {
	if o.IsDeleted() {
		return nil
	}

	now := time.Now()
	o.DeletedAt = &now
	o.Active = false

	return nil
}

func (o *Organization) IsDeleted() bool {
	return o.DeletedAt != nil
}

func (o *Organization) touch() {
	now := time.Now()
	o.UpdatedAt = &now
}

func validateOrganizationName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrOrganizationNameRequired
	}
	return nil
}
