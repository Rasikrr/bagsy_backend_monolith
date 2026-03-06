package workers

import (
	"context"
	"strconv"
	"sync"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/notification"
	"github.com/Rasikrr/core/log"
)

const (
	workersCount = 10
)

type notificationRepository interface {
	PollReady(ctx context.Context, limit int) ([]*notification.Task, error)
	UpdateBatch(ctx context.Context, task []*notification.Task) error
}

type reminderSender interface {
	SendReminder(ctx context.Context, task *notification.Task) error
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

	tasks, err := j.repo.PollReady(ctx, j.batchSize)
	if err != nil {
		log.Error(ctx, "poll notification tasks failed", log.Err(err))
		return
	}
	err = j.processTasks(ctx, tasks)
	if err != nil {
		log.Error(ctx, "process tasks failed", log.Err(err))
		return
	}
	err = j.repo.UpdateBatch(ctx, tasks)
	if err != nil {
		log.Error(ctx, "update tasks failed", log.Err(err))
	}
	log.Infof(ctx, "reminder notification worker finished, processed %d tasks", len(tasks))
}

// nolint: unparam
func (j *ReminderNotificationJob) processTasks(ctx context.Context, tasks []*notification.Task) error {
	taskChan := make(chan *notification.Task, workersCount)
	wg := sync.WaitGroup{}

	wg.Go(func() {
		for {
			select {
			case <-ctx.Done():
				return
			case task, ok := <-taskChan:
				if !ok {
					return
				}
				j.processTask(ctx, task)
			}
		}
	})

	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)
	wg.Wait()
	// Update tasks batch

	return nil
}

func (j *ReminderNotificationJob) processTask(ctx context.Context, task *notification.Task) {
	err := j.sender.SendReminder(ctx, task)
	if err != nil {
		task.MarkFailed(err.Error())
		log.Error(ctx, "send reminder failed",
			log.String("task_id", strconv.FormatInt(task.ID, 10)),
			log.String("type", string(task.Type)),
			log.String("recipient", string(task.RecipientType)),
			log.Err(err),
		)
		return
	}
	task.MarkSent()
}
