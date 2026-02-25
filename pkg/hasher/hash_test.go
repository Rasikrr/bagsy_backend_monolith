package hasher

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

func TestCompare(t *testing.T) {
	password := "Rasik_1234"
	hash := "$2a$10$aKTrWq7Mv7DZB5NZOgklnOvOArT.8M/oRuw9JvsNQG7sW9lYcw/QS"
	out := CheckPassword(hash, password)
	require.True(t, out)
}
