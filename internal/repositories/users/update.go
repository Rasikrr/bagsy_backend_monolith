// nolint: godot
package users

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"time"
)

var (
	errNothingToUpdate = domainErr.NewInvalidInputError("nothing to update", nil)
)

type UserUpdatePatch struct {
	Phone     string
	Name      *string
	Surname   *string
	UpdatedBy *string
}

func (p *UserUpdatePatch) ToSQL() (string, []any, error) {
	if p.IsEmpty() {
		return "", nil, errNothingToUpdate
	}

	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := builder.Update("users").
		Where(sq.And{
			sq.Eq{"phone": p.Phone},
			sq.Eq{"deleted_at": nil},
		})

	if p.Name != nil {
		query = query.Set("name", *p.Name)
	}

	if p.Surname != nil {
		query = query.Set("surname", *p.Surname)
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

func (p *UserUpdatePatch) IsEmpty() bool {
	return p.Name == nil &&
		p.Surname == nil &&
		p.UpdatedBy == nil
}

type UserUpdatePatchBuilder struct {
	patch *UserUpdatePatch
}

func NewUserUpdatePatch() *UserUpdatePatchBuilder {
	return &UserUpdatePatchBuilder{
		patch: &UserUpdatePatch{},
	}
}

func (b *UserUpdatePatchBuilder) SetPhone(phone string) *UserUpdatePatchBuilder {
	b.patch.Phone = phone
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

func (b *UserUpdatePatchBuilder) SetUpdatedBy(updatedBy string) *UserUpdatePatchBuilder {
	b.patch.UpdatedBy = &updatedBy
	return b
}

func (b *UserUpdatePatchBuilder) Build() *UserUpdatePatch {
	return b.patch
}
