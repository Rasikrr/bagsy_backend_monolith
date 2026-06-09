package workers

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/booking"
	"github.com/Rasikrr/core/log"
	"github.com/google/uuid"
)

type appointmentCompletionRepository interface {
	GetDueForCompletion(ctx context.Context, now time.Time) ([]uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*booking.Appointment, error)
	Save(ctx context.Context, a *booking.Appointment) error
}

// AppointmentCompletionJob — воркер авто-завершения прошедших записей
// (confirmed/in_progress с end_at в прошлом → completed). Без него записи
// навсегда застревают в confirmed, и аналитика (выручка/KPI считаются только по
// completed) остаётся пустой.
type AppointmentCompletionJob struct {
	repo     appointmentCompletionRepository
	schedule string
}

func NewAppointmentCompletionJob(repo appointmentCompletionRepository, schedule string) *AppointmentCompletionJob {
	return &AppointmentCompletionJob{
		repo:     repo,
		schedule: schedule,
	}
}

func (j *AppointmentCompletionJob) Name() string {
	return "appointment_completion"
}

func (j *AppointmentCompletionJob) Schedule() string {
	return j.schedule
}

func (j *AppointmentCompletionJob) Run() {
	ctx := context.Background()
	log.Info(ctx, "starting appointment completion worker")

	ids, err := j.repo.GetDueForCompletion(ctx, time.Now())
	if err != nil {
		log.Error(ctx, "failed to get appointments due for completion", log.Err(err))
		return
	}

	var completed int
	for _, id := range ids {
		if j.complete(ctx, id) {
			completed++
		}
	}

	log.Infof(ctx, "appointment completion worker finished, completed %d of %d appointments", completed, len(ids))
}

// complete завершает одну запись. Ошибка одной записи не прерывает остальные.
func (j *AppointmentCompletionJob) complete(ctx context.Context, id uuid.UUID) bool {
	appointment, err := j.repo.GetByID(ctx, id)
	if err != nil {
		log.Error(ctx, "failed to load appointment for completion",
			log.String("appointment_id", id.String()),
			log.Err(err),
		)
		return false
	}

	if err = appointment.AutoComplete(); err != nil {
		log.Error(ctx, "failed to auto-complete appointment",
			log.String("appointment_id", id.String()),
			log.String("status", appointment.Status.String()),
			log.Err(err),
		)
		return false
	}

	if err = j.repo.Save(ctx, appointment); err != nil {
		log.Error(ctx, "failed to save auto-completed appointment",
			log.String("appointment_id", id.String()),
			log.Err(err),
		)
		return false
	}

	return true
}
