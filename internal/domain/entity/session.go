package entity

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/google/uuid"
)

// TODO: remove to other place
type Session struct {
	id          uuid.UUID
	phone       string
	role        enum.Role
	pointCode   string
	networkCode string
}

func NewSession() *Session {
	return &Session{
		id: uuid.New(),
	}
}

func (s *Session) SetPhone(p string) *Session {
	s.phone = p
	return s
}

func (s *Session) SetRole(r enum.Role) *Session {
	s.role = r
	return s
}

func (s *Session) SetPointCode(pc string) *Session {
	s.pointCode = pc
	return s
}

func (s *Session) SetNetworkCode(nc string) *Session {
	s.networkCode = nc
	return s
}

func (s *Session) Phone() string {
	return s.phone
}

func (s *Session) Role() enum.Role {
	return s.role
}

func (s *Session) PointCode() string {
	return s.pointCode
}

func (s *Session) NetworkCode() string {
	return s.networkCode
}

func (s *Session) SessionID() uuid.UUID {
	return s.id
}
