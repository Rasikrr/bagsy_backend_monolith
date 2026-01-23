package bagsynotifications

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/notification"
	"github.com/google/uuid"
)

// notificationsRepository интерфейс для репозитория уведомлений
type notificationsRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*notification.Notification, error)
	GetByBagsyID(ctx context.Context, bagsyID uuid.UUID) ([]*notification.Notification, error)
	GetPendingBatch(ctx context.Context, maxAttempts, limit int) ([]*notification.Notification, error)
	Create(ctx context.Context, n *notification.Notification) (uuid.UUID, error)
	Upsert(ctx context.Context, n *notification.Notification) (uuid.UUID, error)
	CreateBatch(ctx context.Context, notifications []*notification.Notification) error
	MarkSent(ctx context.Context, id uuid.UUID) error
	MarkFailed(ctx context.Context, id uuid.UUID, errMsg string, maxAttempts int) error
	MarkSkipped(ctx context.Context, id uuid.UUID) error
	DeleteByBagsyID(ctx context.Context, bagsyID uuid.UUID) error
	DeleteByBagsyIDs(ctx context.Context, bagsyIDs []uuid.UUID) error
}

type Service struct {
	notificationsRepo notificationsRepository
	maxAttempts       int
}

func NewService(notificationsRepo notificationsRepository, maxAttempts int) *Service {
	return &Service{
		notificationsRepo: notificationsRepo,
		maxAttempts:       maxAttempts,
	}
}

// ScheduleForBagsy создает уведомления для записи
// Реализует интерфейс NotificationScheduler для bagsies service
func (s *Service) ScheduleForBagsy(ctx context.Context, bagsyID uuid.UUID, startAt time.Time) error {
	notifications := buildNotifications(bagsyID, startAt)
	return s.notificationsRepo.CreateBatch(ctx, notifications)
}

// RescheduleForBagsy обновляет уведомления при изменении времени записи
func (s *Service) RescheduleForBagsy(ctx context.Context, bagsyID uuid.UUID, startAt time.Time) error {
	// Upsert обновит существующие или создаст новые
	notifications := buildNotifications(bagsyID, startAt)
	return s.notificationsRepo.CreateBatch(ctx, notifications)
}

// CancelForBagsy удаляет уведомления для записи
func (s *Service) CancelForBagsy(ctx context.Context, bagsyID uuid.UUID) error {
	return s.notificationsRepo.DeleteByBagsyID(ctx, bagsyID)
}

// CancelForBagsies удаляет уведомления для нескольких записей
func (s *Service) CancelForBagsies(ctx context.Context, bagsyIDs []uuid.UUID) error {
	return s.notificationsRepo.DeleteByBagsyIDs(ctx, bagsyIDs)
}

// GetPendingBatch получает batch pending уведомлений для обработки
func (s *Service) GetPendingBatch(ctx context.Context, limit int) ([]*notification.Notification, error) {
	return s.notificationsRepo.GetPendingBatch(ctx, s.maxAttempts, limit)
}

// MarkSent помечает уведомление как отправленное
func (s *Service) MarkSent(ctx context.Context, id uuid.UUID) error {
	return s.notificationsRepo.MarkSent(ctx, id)
}

// MarkFailed помечает уведомление как неудачное
func (s *Service) MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error {
	return s.notificationsRepo.MarkFailed(ctx, id, errMsg, s.maxAttempts)
}

// MarkSkipped помечает уведомление как пропущенное
func (s *Service) MarkSkipped(ctx context.Context, id uuid.UUID) error {
	return s.notificationsRepo.MarkSkipped(ctx, id)
}

// GetByBagsyID получает все уведомления для записи
func (s *Service) GetByBagsyID(ctx context.Context, bagsyID uuid.UUID) ([]*notification.Notification, error) {
	return s.notificationsRepo.GetByBagsyID(ctx, bagsyID)
}

// minTimeBeforeStart минимальное время до записи для отправки уведомления
const minTimeBeforeStart = 5 * time.Minute

// buildNotifications создает список уведомлений для записи
// Логика:
// - Если время уведомления ещё не наступило → планируем на это время
// - Если время уже прошло, но запись ещё не началась → отправляем одно "срочное"
// - Выбираем наиболее подходящее уведомление по времени до записи
func buildNotifications(bagsyID uuid.UUID, startAt time.Time) []*notification.Notification {
	types := notification.AllTypes()
	notifications := make([]*notification.Notification, 0, len(types))

	now := time.Now()
	timeUntilStart := startAt.Sub(now)

	// Если до записи осталось меньше 5 минут — не создаём уведомления
	if timeUntilStart < minTimeBeforeStart {
		return notifications
	}

	recipients := notification.AllRecipients()

	for _, recipient := range recipients {
		// Собираем уведомления: запланированные и одно "срочное" (если нужно) для каждого получателя
		var immediateNotification *notification.Notification

		for _, t := range types {
			scheduledAt := startAt.Add(-time.Duration(t.GetOffset()) * time.Minute)

			if scheduledAt.After(now) {
				// Время ещё не наступило — планируем
				notifications = append(notifications, &notification.Notification{
					BagsyID:       bagsyID,
					Type:          t,
					RecipientType: recipient,
					ScheduledAt:   scheduledAt,
					Status:        notification.StatusPending,
				})
			} else if immediateNotification == nil || t.GetOffset() < immediateNotification.Type.GetOffset() {
				// Время прошло — кандидат на "срочное" уведомление
				// Выбираем наиболее подходящее (с меньшим offset = более актуальное)
				immediateNotification = &notification.Notification{
					BagsyID:       bagsyID,
					Type:          t,
					RecipientType: recipient,
					ScheduledAt:   now,
					Status:        notification.StatusPending,
				}
			}
		}

		// Добавляем "срочное" уведомление, если есть
		if immediateNotification != nil {
			notifications = append(notifications, immediateNotification)
		}
	}

	return notifications
}
