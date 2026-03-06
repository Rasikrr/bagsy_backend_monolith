package notification

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/notification"
	"github.com/cockroachdb/errors"
)

func (u *UseCase) ScheduleReminders(ctx context.Context, params notification.ScheduleParams) error {
	params.Formatter = u.formatter

	tasks := notification.GenerateReminders(params)
	if len(tasks) == 0 {
		return nil
	}

	if err := u.notifRepo.SaveBatch(ctx, tasks); err != nil {
		return errors.Wrap(err, "save reminder tasks")
	}

	return nil
}
