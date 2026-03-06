package notification

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/notification"
	"github.com/google/uuid"
)

type notificationRepository interface {
	SaveBatch(ctx context.Context, tasks []*notification.Task) error
	DeletePendingByAppointmentID(ctx context.Context, appointmentID uuid.UUID) error
}

type UseCase struct {
	notifRepo notificationRepository
	formatter notification.MessageFormatter
}

func NewUseCase(notifRepo notificationRepository, formatter notification.MessageFormatter) *UseCase {
	return &UseCase{
		notifRepo: notifRepo,
		formatter: formatter,
	}
}
