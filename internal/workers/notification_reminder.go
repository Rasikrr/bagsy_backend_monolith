package workers

import (
	"context"
	"strconv"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/notification"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/core/log"
)

type notificationRepository interface {
	PollReady(ctx context.Context, limit int) ([]*notification.Task, error)
	Update(ctx context.Context, task *notification.Task) error
	UnlockStale(ctx context.Context) (int64, error)
}

type reminderSender interface {
	SendReminder(ctx context.Context, phone shared.Phone, message string) error
}

type ReminderNotificationJob struct {
	repo      notificationRepository
	sender    reminderSender
	batchSize int
	schedule  string
}

func NewReminderNotificationJob(
	repo notificationRepository,
	sender reminderSender,
	batchSize int,
	schedule string,
) *ReminderNotificationJob {
	return &ReminderNotificationJob{
		repo:      repo,
		sender:    sender,
		batchSize: batchSize,
		schedule:  schedule,
	}
}

func (j *ReminderNotificationJob) Name() string {
	return "reminder_notification"
}

func (j *ReminderNotificationJob) Schedule() string {
	return j.schedule
}

func (j *ReminderNotificationJob) Run() {
	ctx := context.Background()
	log.Info(ctx, "starting reminder notification worker")

	stale, err := j.repo.UnlockStale(ctx)
	if err != nil {
		log.Error(ctx, "unlock stale notification tasks failed", log.Err(err))
	} else if stale > 0 {
		log.Infof(ctx, "unlocked %d stale notification tasks", stale)
	}

	tasks, err := j.repo.PollReady(ctx, j.batchSize)
	if err != nil {
		log.Error(ctx, "poll notification tasks failed", log.Err(err))
		return
	}

	for _, task := range tasks {
		j.processTask(ctx, task)
	}

	log.Infof(ctx, "reminder notification worker finished, processed %d tasks", len(tasks))
}

func (j *ReminderNotificationJob) processTask(ctx context.Context, task *notification.Task) {
	err := j.sender.SendReminder(ctx, task.RecipientPhone.String(), task.Message)
	if err != nil {
		task.MarkFailed(err.Error())
		log.Error(ctx, "send reminder failed",
			log.String("task_id", strconv.FormatInt(task.ID, 10)),
			log.String("type", string(task.Type)),
			log.String("recipient", string(task.RecipientType)),
			log.Err(err),
		)
	} else {
		task.MarkSent()
	}

	if updateErr := j.repo.Update(ctx, task); updateErr != nil {
		log.Error(ctx, "update notification task failed",
			log.String("task_id", strconv.FormatInt(task.ID, 10)),
			log.Err(updateErr),
		)
	}
}
