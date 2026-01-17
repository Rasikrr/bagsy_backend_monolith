package actor

import (
	"context"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	"github.com/google/uuid"
)

type Actor struct {
	id          uuid.UUID
	phone       string
	role        user.Role
	pointCode   string
	networkCode string
}

func NewActor() *Actor {
	return &Actor{
		id: uuid.New(),
	}
}

func (s *Actor) SetPhone(p string) *Actor {
	s.phone = p
	return s
}

func (s *Actor) SetRole(r user.Role) *Actor {
	s.role = r
	return s
}

func (s *Actor) SetRoleString(r string) (*Actor, error) {
	role, err := user.RoleString(r)
	if err != nil {
		return nil, ErrUnknownRole.WithError(err)
	}
	s.SetRole(role)
	return s, nil
}

func (s *Actor) SetPointCode(pc string) *Actor {
	s.pointCode = pc
	return s
}

func (s *Actor) SetNetworkCode(nc string) *Actor {
	s.networkCode = nc
	return s
}

func (s *Actor) Phone() string {
	return s.phone
}

func (s *Actor) Role() user.Role {
	return s.role
}

func (s *Actor) PointCode() string {
	return s.pointCode
}

func (s *Actor) NetworkCode() string {
	return s.networkCode
}

func (s *Actor) ID() uuid.UUID {
	return s.id
}

type sessionKey struct{}

func GetActor(ctx context.Context) (*Actor, error) {
	if session, ok := ctx.Value(sessionKey{}).(*Actor); ok {
		return session, nil
	}
	return nil, domainErr.NewUnauthorizedError("actor not found")
}

func SetActor(ctx context.Context, session *Actor) context.Context {
	return context.WithValue(ctx, sessionKey{}, session)
}
