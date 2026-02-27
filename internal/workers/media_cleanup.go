package workers

import (
	"context"
	"time"

	"github.com/Rasikrr/core/log"
)

type mediaRepository interface {
	MarkExpiredPendingAsFailed(ctx context.Context, threshold time.Time) (int64, error)
}

type MediaCleanupJob struct {
	repo      mediaRepository
	threshold time.Duration
	schedule  string
}

func NewMediaCleanupJob(repo mediaRepository, threshold time.Duration, schedule string) *MediaCleanupJob {
	return &MediaCleanupJob{
		repo:      repo,
		threshold: threshold,
		schedule:  schedule,
	}
}

func (j *MediaCleanupJob) Name() string {
	return "media_cleanup"
}

func (j *MediaCleanupJob) Schedule() string {
	return j.schedule
}

func (j *MediaCleanupJob) Run() {
	ctx := context.Background()
	log.Info(ctx, "starting media cleanup worker")

	expirationTime := time.Now().Add(-j.threshold)
	affected, err := j.repo.MarkExpiredPendingAsFailed(ctx, expirationTime)
	if err != nil {
		log.Error(ctx, "media cleanup worker failed", log.Err(err))
		return
	}

	log.Infof(ctx, "media cleanup worker finished, marked %d assets as failed", affected)
}
