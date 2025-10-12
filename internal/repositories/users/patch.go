// nolint: godot
package users

import (
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
)

var (
	builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type UserUpdatePatch struct {
	Phones    []string
	Role      *enum.Role
	Name      *string
	Surname   *string
	PointCode *string
	Network   *string
	Active    *bool
	Password  *string
	UpdatedBy *string

	// UpdatedAt всегда проставляется автоматически при любом обновлении
	// Не включаем в структуру, будет задаваться в репозитории
}

func (p *UserUpdatePatch) ToSQL() (string, []any, error) {
	if p.IsEmpty() {
		return "", nil, errNothingToUpdate
	}

	query := builder.Update("users").
		Where(squirrel.Eq{"phone": p.Phones})

	if p.Role != nil {
		query = query.Set("role", p.Role.String())
	}

	if p.Name != nil {
		query = query.Set("name", *p.Name)
	}

	if p.Surname != nil {
		query = query.Set("surname", *p.Surname)
	}

	if p.PointCode != nil {
		if *p.PointCode == "" {
			query = query.Set("point_code", nil)
		} else {
			query = query.Set("point_code", *p.PointCode)
		}
	}

	if p.Network != nil {
		if *p.Network == "" {
			query = query.Set("network", nil)
		} else {
			query = query.Set("network", *p.Network)
		}
	}

	if p.Active != nil {
		query = query.Set("active", *p.Active)
	}

	if p.Password != nil {
		if *p.Password == "" {
			query = query.Set("password", nil)
		} else {
			query = query.Set("password", *p.Password)
		}
	}

	if p.UpdatedBy != nil {
		query = query.Set("updated_by", *p.UpdatedBy)
	}

	query = query.Set("updated_at", time.Now())

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("build update query: %w", err)
	}

	return sql, args, nil
}

type UserUpdatePatchBuilder struct {
	patch *UserUpdatePatch
}

func NewUserUpdatePatch() *UserUpdatePatchBuilder {
	return &UserUpdatePatchBuilder{
		patch: &UserUpdatePatch{},
	}
}

func (b *UserUpdatePatchBuilder) SetPhones(phones ...string) *UserUpdatePatchBuilder {
	b.patch.Phones = phones
	return b
}

func (b *UserUpdatePatchBuilder) SetRole(role enum.Role) *UserUpdatePatchBuilder {
	b.patch.Role = &role
	return b
}

func (b *UserUpdatePatchBuilder) SetName(name string) *UserUpdatePatchBuilder {
	b.patch.Name = &name
	return b
}

func (b *UserUpdatePatchBuilder) SetSurname(surname string) *UserUpdatePatchBuilder {
	b.patch.Surname = &surname
	return b
}

func (b *UserUpdatePatchBuilder) SetPointCode(pointCode string) *UserUpdatePatchBuilder {
	b.patch.PointCode = &pointCode
	return b
}

func (b *UserUpdatePatchBuilder) ClearPointCode() *UserUpdatePatchBuilder {
	empty := ""
	b.patch.PointCode = &empty
	return b
}

func (b *UserUpdatePatchBuilder) SetActive(active bool) *UserUpdatePatchBuilder {
	b.patch.Active = &active
	return b
}

func (b *UserUpdatePatchBuilder) SetPassword(password string) *UserUpdatePatchBuilder {
	b.patch.Password = &password
	return b
}

func (b *UserUpdatePatchBuilder) ClearPassword() *UserUpdatePatchBuilder {
	empty := ""
	b.patch.Password = &empty
	return b
}

func (b *UserUpdatePatchBuilder) SetUpdatedBy(updatedBy string) *UserUpdatePatchBuilder {
	b.patch.UpdatedBy = &updatedBy
	return b
}

func (b *UserUpdatePatchBuilder) Build() *UserUpdatePatch {
	return b.patch
}

func (p *UserUpdatePatch) IsEmpty() bool {
	return p.Role == nil &&
		p.Name == nil &&
		p.Surname == nil &&
		p.PointCode == nil &&
		p.Active == nil &&
		p.Password == nil &&
		p.UpdatedBy == nil
}

// GetUpdatedFields возвращает map полей для обновления (для логирования/отладки)
func (p *UserUpdatePatch) GetUpdatedFields() map[string]interface{} {
	fields := make(map[string]interface{})

	if p.Role != nil {
		fields["role"] = *p.Role
	}
	if p.Name != nil {
		fields["name"] = *p.Name
	}
	if p.Surname != nil {
		fields["surname"] = *p.Surname
	}
	if p.PointCode != nil {
		if *p.PointCode == "" {
			fields["point_code"] = nil // NULL
		} else {
			fields["point_code"] = *p.PointCode
		}
	}
	if p.Active != nil {
		fields["active"] = *p.Active
	}
	if p.Password != nil {
		if *p.Password == "" {
			fields["password"] = nil // NULL
		} else {
			fields["password"] = "[REDACTED]" // не логируем пароли
		}
	}
	if p.UpdatedBy != nil {
		fields["updated_by"] = *p.UpdatedBy
	}

	// UpdatedAt всегда обновляется
	fields["updated_at"] = time.Now()

	return fields
}
