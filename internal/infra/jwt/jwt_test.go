package jwt

import (
	"testing"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	testSecret = "bagsy123"
	testIssuer = "bagsies-test"
)

func TestAccessToken(t *testing.T) {
	mgr := NewTokenManager(testSecret, testIssuer)
	ttl := time.Minute * 15

	expiresAt := time.Now().Add(ttl)
	userID := uuid.New()
	phone, err := shared.NewPhone("+77715275251")
	require.NoError(t, err)

	authToken := auth.NewToken(userID, phone, expiresAt)

	token, err := mgr.NewAccessToken(authToken)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	t.Log(token)
}

func TestRefreshToken(t *testing.T) {
	mgr := NewTokenManager(testSecret, testIssuer)
	refresh, hash, err := mgr.NewRefreshToken()
	require.NoError(t, err)
	require.NotEmpty(t, refresh)
	require.NotEmpty(t, hash)
	t.Logf("refresh: %s, hash: %s", refresh, hash)
}
