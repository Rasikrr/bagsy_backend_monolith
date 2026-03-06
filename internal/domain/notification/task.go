package notification

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

type Metadata struct {
	ServiceName   string
	LocationName  string
	AppointmentAt time.Time
}

type Task struct {
	ID             int64
	AppointmentID  uuid.UUID
	Type           Type
	RecipientType  RecipientType
	RecipientPhone shared.Phone
	Metadata       Metadata
	Status         Status
	ScheduledFor   time.Time
	Attempts       int
	MaxAttempts    int
	LastError      *string
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}

type CreateTaskParams struct {
	AppointmentID  uuid.UUID
	Type           Type
	RecipientType  RecipientType
	RecipientPhone shared.Phone
	Metadata       Metadata
	ScheduledFor   time.Time
	MaxAttempts    int
}

func NewTask(params CreateTaskParams) *Task {
	return &Task{
		AppointmentID:  params.AppointmentID,
		Type:           params.Type,
		RecipientType:  params.RecipientType,
		RecipientPhone: params.RecipientPhone,
		Metadata:       params.Metadata,
		Status:         StatusPending,
		ScheduledFor:   params.ScheduledFor,
		Attempts:       0,
		MaxAttempts:    params.MaxAttempts,
		CreatedAt:      time.Now(),
	}
}

// MarkSent transitions to sent after successful delivery.
func (t *Task) MarkSent() {
	t.Status = StatusSent
	t.Attempts++
	t.touch()
}

// MarkFailed increments attempts and sets error. Returns to pending if retries remain, otherwise failed.
func (t *Task) MarkFailed(errMsg string) {
	t.Attempts++
	t.LastError = &errMsg
	if t.Attempts >= t.MaxAttempts {
		t.Status = StatusFailed
	} else {
		t.Status = StatusPending
	}
	t.touch()
}

// CanRetry returns true if the task has remaining attempts.
func (t *Task) CanRetry() bool {
	return t.Attempts < t.MaxAttempts
}

func (t *Task) touch() {
	now := time.Now()
	t.UpdatedAt = &now
}
