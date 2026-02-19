package access

import (
	"context"
	"testing"

	"github.com/Rasikrr/core/database/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_AccessRepo(t *testing.T) {
	ctx := context.Background()
	userID := uuid.MustParse("2e642a3e-4223-4473-97a0-367691dcfdbc")
	db, err := postgres.NewPostgres(ctx, postgres.Config{
		DSN:                 "postgresql://postgres:rasik1234@localhost:5432/bagsy",
		Required:            true,
		MaxConns:            1,
		MinConns:            1,
		MaxIdleConnIdleTime: 1,
	})
	require.NoError(t, err)
	repo := NewRepository(db)
	out, err := repo.GetOrgContext(ctx, userID)
	require.NoError(t, err)
	t.Logf("out: %+v", out)
}
