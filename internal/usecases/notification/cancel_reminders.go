package notification

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

func (u *UseCase) CancelReminders(ctx context.Context, appointmentID uuid.UUID) error {
	if err := u.notifRepo.DeletePendingByAppointmentID(ctx, appointmentID); err != nil {
		return errors.Wrap(err, "cancel reminders")
	}

	return nil
}
