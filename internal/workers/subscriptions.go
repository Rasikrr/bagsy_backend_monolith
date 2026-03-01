package workers

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/core/log"
)

// TODO: review and refactor
const (
	// suspensionDays — количество дней в suspended до автоматической отмены.
	suspensionDays = 90
)

type subscriptionRepository interface {
	GetRequiringAction(ctx context.Context, now time.Time) ([]*billing.Subscription, error)
	GetPendingDeletion(ctx context.Context, now time.Time) ([]*billing.Subscription, error)
	Save(ctx context.Context, sub *billing.Subscription) error
}

type SubscriptionStatusJob struct {
	repo     subscriptionRepository
	schedule string
}

func NewSubscriptionStatusJob(repo subscriptionRepository, schedule string) *SubscriptionStatusJob {
	return &SubscriptionStatusJob{
		repo:     repo,
		schedule: schedule,
	}
}

func (j *SubscriptionStatusJob) Name() string {
	return "subscription_status"
}

func (j *SubscriptionStatusJob) Schedule() string {
	return j.schedule
}

func (j *SubscriptionStatusJob) Run() {
	ctx := context.Background()
	log.Info(ctx, "starting subscription status worker")

	now := time.Now()

	j.processStatusTransitions(ctx, now)
	j.processDataDeletion(ctx, now)

	log.Info(ctx, "subscription status worker finished")
}

func (j *SubscriptionStatusJob) processStatusTransitions(ctx context.Context, now time.Time) {
	subs, err := j.repo.GetRequiringAction(ctx, now)
	if err != nil {
		log.Error(ctx, "failed to get subscriptions requiring action", log.Err(err))
		return
	}

	for _, sub := range subs {
		j.transitionSubscription(ctx, sub)
	}
}

func (j *SubscriptionStatusJob) transitionSubscription(ctx context.Context, sub *billing.Subscription) {
	var (
		err        error
		transition string
	)

	switch {
	case sub.ShouldSuspend():
		transition = "suspend"
		err = sub.Suspend()

	case sub.NeedsRetry():
		transition = "schedule_retry"
		err = sub.ScheduleRetry()

	case sub.ShouldCancel(suspensionDays):
		transition = "cancel"
		err = sub.Cancel()

	case sub.IsPendingCancellation():
		transition = "cancel_at_period_end"
		err = sub.Cancel()

	case sub.Status == billing.SubscriptionStatusTrial || sub.Status == billing.SubscriptionStatusActive:
		transition = "mark_past_due"
		err = sub.MarkPastDue()

	default:
		return
	}

	if err != nil {
		log.Error(ctx, "subscription transition failed",
			log.String("subscription_id", sub.ID.String()),
			log.String("org_id", sub.OrganizationID.String()),
			log.String("transition", transition),
			log.Err(err),
		)
		return
	}

	if err = j.repo.Save(ctx, sub); err != nil {
		log.Error(ctx, "failed to save subscription after transition",
			log.String("subscription_id", sub.ID.String()),
			log.String("transition", transition),
			log.Err(err),
		)
		return
	}

	log.Infof(ctx, "subscription %s transitioned to %s (org: %s)",
		sub.ID, sub.Status, sub.OrganizationID)
}

func (j *SubscriptionStatusJob) processDataDeletion(ctx context.Context, now time.Time) {
	subs, err := j.repo.GetPendingDeletion(ctx, now)
	if err != nil {
		log.Error(ctx, "failed to get subscriptions pending deletion", log.Err(err))
		return
	}

	for _, sub := range subs {
		// TODO: реализовать удаление данных организации (employees, locations, services, etc.)
		log.Infof(ctx, "subscription %s (org: %s) is due for data deletion — not yet implemented",
			sub.ID, sub.OrganizationID)
	}
}
