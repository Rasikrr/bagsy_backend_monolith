package jwt

import (
	"testing"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	authS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"
	"github.com/stretchr/testify/require"
)

const (
	testSecret = "bagsy123"
	testIssuer = "bagsies-test"
)

func TestAccessToken(t *testing.T) {
	mgr := NewTokenManager(testSecret, testIssuer)
	ttl := time.Minute * 15
	payload := &authS.AccessTokenPayload{
		Phone:       "77715275251",
		Role:        user.RoleStaff.String(),
		PointCode:   "test_point",
		NetworkCode: "test_network",
	}
	token, err := mgr.NewAccessToken(payload, ttl)
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
