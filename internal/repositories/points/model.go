package points

import (
	"encoding/json"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/hashicorp/go-multierror"
)

type model struct {
	Code        string     `db:"code"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	NetworkCode string     `db:"network_code"`
	CategoryID  int        `db:"category_id"`
	Address     []byte     `db:"address"`
	City        string     `db:"city"`
	Active      bool       `db:"active"`
	Schedule    []byte     `db:"schedule"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	UpdatedBy   string     `db:"updated_by"`
}

type models []model

func convert(e *entity.Point) (model, error) {
	var mErr *multierror.Error
	var m model
	m.Code = e.Code
	m.Name = e.Name
	m.Description = e.Description
	m.NetworkCode = e.NetworkCode
	m.CategoryID = e.CategoryID
	m.City = e.City
	m.Active = e.Active
	m.CreatedAt = e.CreatedAt
	m.UpdatedAt = e.UpdatedAt
	m.DeletedAt = e.DeletedAt
	m.UpdatedBy = e.UpdatedBy
	var err error
	if m.Address, err = json.Marshal(e.Address); err != nil {
		mErr = multierror.Append(mErr, err)
	}
	if m.Schedule, err = json.Marshal(e.Schedule); err != nil {
		mErr = multierror.Append(mErr, err)
	}
	return m, mErr.ErrorOrNil()
}

func (m model) convert() *entity.Point {
	point := &entity.Point{
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		NetworkCode: m.NetworkCode,
		CategoryID:  m.CategoryID,
		City:        m.City,
		Active:      m.Active,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
		UpdatedBy:   m.UpdatedBy,
	}

	json.Unmarshal(m.Address, &point.Address)
	json.Unmarshal(m.Schedule, &point.Schedule)

	return point
}
