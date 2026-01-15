package points

import (
	"encoding/json"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/samber/lo"
)

/*
code         TEXT PRIMARY KEY,
name         TEXT NOT NULL,
description  TEXT,
network_code TEXT NOT NULL,
category_id  INTEGER NOT NULL,
address      JSONB NOT NULL DEFAULT '{}',
city         TEXT NOT NULL,
active       BOOLEAN NOT NULL DEFAULT false,
schedule     JSONB NOT NULL DEFAULT '{}',
created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at   TIMESTAMPTZ DEFAULT now(),
deleted_at   TIMESTAMPTZ,
updated_by   TEXT NOT NULL DEFAULT 'system',
CONSTRAINT point_code_network_code_unique UNIQUE (code, network_code)
*/
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

func (mm models) convert() []*point.Point {
	return lo.Map(mm, func(item model, _ int) *point.Point {
		return item.convert()
	})
}

func convert(e *point.Point) (model, error) {
	out := model{
		Code:        e.Code,
		Name:        e.Name,
		Description: e.Description,
		NetworkCode: e.NetworkCode,
		CategoryID:  e.CategoryID,
		City:        e.City,
		Active:      e.Active,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		DeletedAt:   e.DeletedAt,
		UpdatedBy:   e.UpdatedBy,
	}

	addressDto := addressToDTO(e.Address)
	bb, err := json.Marshal(addressDto)
	if err != nil {
		return out, err
	}
	out.Address = bb

	schedulesDto := schedulesToDTO(e.Schedule)
	bb, err = json.Marshal(schedulesDto)
	if err != nil {
		return out, err
	}
	out.Schedule = bb
	return out, nil
}

func (m model) convert() *point.Point {
	point := &point.Point{
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

	var addrDTO addressDTO
	if err := json.Unmarshal(m.Address, &addrDTO); err == nil {
		point.Address = addrDTO.toEntity()
	}

	var scheduleDTOs schedulesDTO
	if err := json.Unmarshal(m.Schedule, &scheduleDTOs); err == nil {
		point.Schedule = scheduleDTOs.toEntity()
	}

	return point
}
