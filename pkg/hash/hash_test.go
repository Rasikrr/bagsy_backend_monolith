package hash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	password := "admin123"
	hashedPassword, err := Password(password)
	t.Logf("Hashed password: %s", hashedPassword)
	require.NoError(t, err)
}
