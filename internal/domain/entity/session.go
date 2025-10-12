package entity

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/google/uuid"
)

type Session struct {
	ID          uuid.UUID
	Phone       string
	Active      bool
	Role        enum.Role
	PointCode   *string
	NetworkCode *string
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
	s.PointCode = &pc
	return s
}

func (s *Session) SetNetworkCode(nc string) *Session {
	s.NetworkCode = &nc
	return s
}

func (s *Session) GetPhone() string {
	return s.Phone
}

func (s *Session) GetRole() enum.Role {
	return s.Role
}

func (s *Session) GetActive() bool {
	return s.Active
}

func (s *Session) GetPointCode() string {
	if s.PointCode != nil {
		return *s.PointCode
	}
	return ""
}

func (s *Session) GetNetworkCode() string {
	if s.NetworkCode != nil {
		return *s.NetworkCode
	}
	return ""
}

func (s *Session) GetID() uuid.UUID {
	return s.ID
}
