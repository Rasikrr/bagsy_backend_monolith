package session

import (
	"context"
	"errors"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
)

type ContextKey string

const (
	SessionKey ContextKey = "_session"
)

func GetSession(ctx context.Context) (*entity.Session, error) {
	if session, ok := ctx.Value(SessionKey).(*entity.Session); ok {
		return session, nil
	}
	return nil, errors.New("session not found")
}

func SetSession(ctx context.Context, session *entity.Session) context.Context {
	return context.WithValue(ctx, SessionKey, session)
}
