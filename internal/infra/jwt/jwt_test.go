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
	testSecret = "JWT_SUPER_SECRET"
	testIssuer = "TEST_ENVIROMENT_ERNAR"
)

func TestAccessToken(t *testing.T) {
	mgr := NewTokenManager(testSecret, testIssuer)
	ttl := time.Minute * 60 * 24 * 7

	userID := uuid.New()
	phone, err := shared.NewPhone("+77715275251")
	require.NoError(t, err)

	authToken := auth.NewToken(userID, phone, ttl)

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
