package access

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
)

type ctxKey int

const (
	tokenKey ctxKey = iota
	orgContextKey
)

func WithToken(ctx context.Context, t *auth.Token) context.Context {
	return context.WithValue(ctx, tokenKey, t)
}

func TokenFromContext(ctx context.Context) (*auth.Token, bool) {
	t, ok := ctx.Value(tokenKey).(*auth.Token)
	return t, ok
}

func WithOrgContext(ctx context.Context, oc *OrgContext) context.Context {
	return context.WithValue(ctx, orgContextKey, oc)
}

func OrgContextFromContext(ctx context.Context) (*OrgContext, bool) {
	oc, ok := ctx.Value(orgContextKey).(*OrgContext)
	return oc, ok
}
