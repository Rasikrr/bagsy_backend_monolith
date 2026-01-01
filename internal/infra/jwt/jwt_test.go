package jwt

import (
	"testing"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/util"
	"github.com/stretchr/testify/require"
)

const (
	testSecret = "bagsy123"
	testIssuer = "bagsy-test"
)

func TestAccessToken(t *testing.T) {
	mgr := NewTokenManager(testSecret, testIssuer)
	ttl := time.Minute * 15
	user := &entity.User{
		Phone:       "77715275251",
		Role:        enum.RoleStaff,
		Name:        "Rassul",
		Surname:     "Turtulov",
		PointCode:   util.Pointer("test_point"),
		NetworkCode: util.Pointer("test_network"),
		Active:      true,
	}
	token, err := mgr.NewAccessToken(user, ttl)
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
