package entity

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID
	Phone     string
	Active    bool
	Role      enum.Role
	PointCode string
}

func NewSession() *Session {
	return &Session{
		ID: uuid.New(),
	}
}

func (s *Session) SetPhone(p string) *Session {
	s.Phone = p
	return s
}

func (s *Session) SetRole(r enum.Role) *Session {
	s.Role = r
	return s
}

func (s *Session) SetActive(a bool) *Session {
	s.Active = a
	return s
}

func (s *Session) SetPointCode(pc string) *Session {
	s.PointCode = pc
	return s
}
