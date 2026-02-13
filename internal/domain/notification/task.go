package notification

import (
	"time"

	"github.com/google/uuid"
)

// Task represents a record in notification_outbox for Transactional Outbox pattern.
type Task struct {
	ID           int64
	EntityID     string
	Type         Type
	Payload      Payload
	ScheduledFor time.Time
	LockedUntil  *time.Time
	CreatedAt    time.Time
}

type Payload struct {
	CustomerID uuid.UUID
	MasterID   uuid.UUID
	Message    string
}

func NewPayload(customerID, masterID uuid.UUID, msg string) Payload {
	return Payload{
		CustomerID: customerID,
		MasterID:   masterID,
		Message:    msg,
	}
}

func NewNotificationTask(
	entityID string,
	taskType Type,
	payload Payload,
	scheduledFor time.Time,
) (*Task, error) {
	return &Task{
		EntityID:     entityID,
		Type:         taskType,
		Payload:      payload,
		ScheduledFor: scheduledFor,
		CreatedAt:    time.Now(),
	}, nil
}

func (t *Task) Lock(duration time.Duration) {
	until := time.Now().Add(duration)
	t.LockedUntil = &until
}

func (t *Task) IsLocked() bool {
	if t.LockedUntil == nil {
		return false
	}
	return t.LockedUntil.After(time.Now())
}
