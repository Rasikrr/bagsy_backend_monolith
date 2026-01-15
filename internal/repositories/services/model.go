package services

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/google/uuid"
)

/*
id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
point_code       TEXT NOT NULL,
category_id      INTEGER NOT NULL,
subcategory_id   INTEGER,
name             TEXT NOT NULL,
description      TEXT,
duration_minutes INTEGER NOT NULL,
active           BOOLEAN DEFAULT false,
created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at       TIMESTAMPTZ DEFAULT now(),
updated_by       TEXT NOT NULL DEFAULT 'system'
*/
type model struct {
	ID              uuid.UUID  `db:"id"`
	PointCode       string     `db:"point_code"`
	CategoryID      int        `db:"category_id"`
	SubcategoryID   *int       `db:"subcategory_id"`
	Name            string     `db:"name"`
	Description     *string    `db:"description"`
	DurationMinutes int        `db:"duration_minutes"`
	Active          bool       `db:"active"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
	UpdatedBy       *string    `db:"updated_by"`
}

func convert(e *service.Service) model {
	return model{
		ID:              e.ID,
		PointCode:       e.PointCode,
		CategoryID:      e.CategoryID,
		SubcategoryID:   e.SubcategoryID,
		Name:            e.Name,
		Description:     e.Description,
		DurationMinutes: e.DurationMinutes,
		Active:          e.Active,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
		UpdatedBy:       e.UpdatedBy,
	}
}

func (m model) convert() *service.Service {
	return &service.Service{
		ID:              m.ID,
		PointCode:       m.PointCode,
		CategoryID:      m.CategoryID,
		SubcategoryID:   m.SubcategoryID,
		Name:            m.Name,
		Description:     m.Description,
		DurationMinutes: m.DurationMinutes,
		Active:          m.Active,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		UpdatedBy:       m.UpdatedBy,
	}
}

type models []model

func (m models) convert() []*service.Service {
	list := make([]*service.Service, len(m))
	for i, item := range m {
		list[i] = item.convert()
	}
	return list
}
