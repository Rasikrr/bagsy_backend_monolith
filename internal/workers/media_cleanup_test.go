package workers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type mockMediaRepository struct {
	mock.Mock
}

func (m *mockMediaRepository) MarkExpiredPendingAsFailed(ctx context.Context, threshold time.Time) (int64, error) {
	args := m.Called(ctx, threshold)
	return int64(args.Int(0)), args.Error(1)
}

func TestMediaCleanupJob_Run(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mockMediaRepository)
		threshold := 15 * time.Minute
		job := NewMediaCleanupJob(repo, threshold, "0 */15 * * * *")

		repo.On("MarkExpiredPendingAsFailed", mock.Anything, mock.AnythingOfType("time.Time")).
			Return(5, nil)

		job.Run()

		repo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		repo := new(mockMediaRepository)
		threshold := 15 * time.Minute
		job := NewMediaCleanupJob(repo, threshold, "0 */15 * * * *")

		repo.On("MarkExpiredPendingAsFailed", mock.Anything, mock.AnythingOfType("time.Time")).
			Return(0, errors.New("db error"))

		job.Run()

		repo.AssertExpectations(t)
	})
}
