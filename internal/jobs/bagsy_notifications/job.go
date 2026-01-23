package bagsynotifications

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/bagsy"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/notification"
	"github.com/Rasikrr/core/log"
	"github.com/google/uuid"
)

// notificationService интерфейс для сервиса уведомлений (расписание)
type notificationService interface {
	GetPendingBatch(ctx context.Context, limit int) ([]*notification.Notification, error)
	MarkSent(ctx context.Context, id uuid.UUID) error
	MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error
	MarkSkipped(ctx context.Context, id uuid.UUID) error
}

// bagsyService интерфейс для получения информации о записи
type bagsyService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*bagsy.Bagsy, error)
}

// messenger интерфейс для отправки сообщений (SMS/WhatsApp)
type messenger interface {
	SendBagsyReminder(ctx context.Context, phone string, b *bagsy.Bagsy, notifType notification.Type, recipientType notification.RecipientType) error
}

// processResult результат обработки уведомления
type processResult int

const (
	resultSuccess processResult = iota
	resultFailed
	resultSkipped
)

// Job обрабатывает отправку уведомлений о записях с использованием worker pool
type Job struct {
	name                string
	schedule            string
	batchSize           int
	workerCount         int
	notificationService notificationService
	bagsyService        bagsyService
	messenger           messenger
}

// NewJob создает новую джобу для обработки уведомлений
// schedule: cron расписание, например "0 */1 * * * *" (каждую минуту с секундами)
// workerCount: количество параллельных воркеров
func NewJob(
	name string,
	schedule string,
	batchSize int,
	workerCount int,
	notificationService notificationService,
	bagsyService bagsyService,
	messenger messenger,
) *Job {
	if workerCount < 1 {
		workerCount = 1
	}
	return &Job{
		name:                name,
		schedule:            schedule,
		batchSize:           batchSize,
		workerCount:         workerCount,
		notificationService: notificationService,
		bagsyService:        bagsyService,
		messenger:           messenger,
	}
}

// Name возвращает имя джобы
func (j *Job) Name() string {
	return j.name
}

// Schedule возвращает cron расписание
func (j *Job) Schedule() string {
	return j.schedule
}

// Run выполняет обработку pending уведомлений с использованием worker pool
func (j *Job) Run() {
	ctx := context.Background()
	log.Infof(ctx, "[%s] Starting notification processing (workers=%d)", j.name, j.workerCount)

	notifications, err := j.notificationService.GetPendingBatch(ctx, j.batchSize)
	if err != nil {
		log.Errorf(ctx, "[%s] Failed to get pending notifications: %v", j.name, err)
		return
	}

	if len(notifications) == 0 {
		log.Debugf(ctx, "[%s] No pending notifications to process", j.name)
		return
	}

	log.Infof(ctx, "[%s] Processing %d notifications with %d workers", j.name, len(notifications), j.workerCount)

	// Счетчики результатов (atomic для thread-safety)
	var successCount, failCount, skipCount int64

	// Канал задач
	jobs := make(chan *notification.Notification, len(notifications))

	// WaitGroup для ожидания завершения всех воркеров
	var wg sync.WaitGroup

	// Запускаем воркеры
	for w := range j.workerCount {
		wg.Add(1)
		go j.worker(ctx, w, jobs, &wg, &successCount, &failCount, &skipCount)
	}

	// Отправляем задачи в канал
	for _, n := range notifications {
		jobs <- n
	}
	close(jobs) // Закрываем канал после отправки всех задач

	// Ждем завершения всех воркеров
	wg.Wait()

	log.Infof(ctx, "[%s] Finished: success=%d, failed=%d, skipped=%d",
		j.name, successCount, failCount, skipCount)
}

// worker обрабатывает уведомления из канала
func (j *Job) worker(
	ctx context.Context,
	id int,
	jobs <-chan *notification.Notification,
	wg *sync.WaitGroup,
	successCount, failCount, skipCount *int64,
) {
	defer wg.Done()

	for n := range jobs {
		result := j.processNotification(ctx, n)

		switch result {
		case resultSuccess:
			atomic.AddInt64(successCount, 1)
		case resultFailed:
			atomic.AddInt64(failCount, 1)
		case resultSkipped:
			atomic.AddInt64(skipCount, 1)
		}
	}

	log.Debugf(ctx, "[%s] Worker %d finished", j.name, id)
}

func (j *Job) processNotification(ctx context.Context, n *notification.Notification) processResult {
	// Получаем информацию о записи
	b, err := j.bagsyService.GetByID(ctx, n.BagsyID)
	if err != nil {
		log.Warnf(ctx, "[%s] Failed to get bagsy %s: %v", j.name, n.BagsyID, err)
		// Если запись не найдена — пропускаем уведомление
		if markErr := j.notificationService.MarkSkipped(ctx, n.ID); markErr != nil {
			log.Errorf(ctx, "[%s] Failed to mark notification %s as skipped: %v", j.name, n.ID, markErr)
		}
		return resultSkipped
	}

	// Проверяем, не прошло ли уже время записи
	if b.StartAt.Before(time.Now()) {
		log.Infof(ctx, "[%s] Bagsy %s start time has passed, skipping notification", j.name, n.BagsyID)
		if markErr := j.notificationService.MarkSkipped(ctx, n.ID); markErr != nil {
			log.Errorf(ctx, "[%s] Failed to mark notification %s as skipped: %v", j.name, n.ID, markErr)
		}
		return resultSkipped
	}

	// Определяем номер телефона получателя
	var phone string
	switch n.RecipientType {
	case notification.RecipientTypeClient:
		phone = b.ClientPhone
	case notification.RecipientTypeMaster:
		phone = b.MasterPhone
	default:
		log.Errorf(ctx, "[%s] Unknown recipient type %s for notification %s", j.name, n.RecipientType, n.ID)
		if markErr := j.notificationService.MarkFailed(ctx, n.ID, "unknown recipient type"); markErr != nil {
			log.Errorf(ctx, "[%s] Failed to mark notification %s as failed: %v", j.name, n.ID, markErr)
		}
		return resultFailed
	}

	// Отправляем уведомление
	err = j.messenger.SendBagsyReminder(ctx, phone, b, n.Type, n.RecipientType)
	if err != nil {
		log.Errorf(ctx, "[%s] Failed to send notification %s: %v", j.name, n.ID, err)
		if markErr := j.notificationService.MarkFailed(ctx, n.ID, err.Error()); markErr != nil {
			log.Errorf(ctx, "[%s] Failed to mark notification %s as failed: %v", j.name, n.ID, markErr)
		}
		return resultFailed
	}

	// Помечаем как отправленное
	if markErr := j.notificationService.MarkSent(ctx, n.ID); markErr != nil {
		log.Errorf(ctx, "[%s] Failed to mark notification %s as sent: %v", j.name, n.ID, markErr)
		return resultFailed
	}

	log.Infof(ctx, "[%s] Successfully sent notification %s (type=%s, recipient=%s) for bagsy %s",
		j.name, n.ID, n.Type.String(), n.RecipientType.String(), n.BagsyID)

	return resultSuccess
}
