package users

import (
	"context"
	"testing"
	"time"

	"github.com/Rasikrr/core/config"
	"github.com/Rasikrr/core/database"
	"github.com/stretchr/testify/require"
)

func TestRepository_Update(t *testing.T) {
	t.Skip()
	ctx := context.Background()
	db, err := database.NewPostgres(ctx, config.PostgresConfig{
		DSN:                 "postgres://postgres:rasik1234@localhost:5432/bugsy?sslmode=disable",
		Required:            true,
		MaxConns:            10,
		MinConns:            3,
		MaxIdleConnIdleTime: 1 * time.Minute,
	})
	repo := NewRepository(db)
	patch := NewUserUpdatePatch()
	patch.SetPhones("77777777777", "77715275253", "2133234243232")
	patch.SetPassword("123456")
	patch.SetActive(true)
	patch.SetUpdatedBy("test")

	out := patch.Build()

	sq, args, err := out.ToSQL()
	require.NoError(t, err)
	t.Logf("sql: %s, args: %v", sq, args)

	err = repo.Update(ctx, out)
	require.NoError(t, err)
}
