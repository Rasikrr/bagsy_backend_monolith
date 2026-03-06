package notification

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

const (
	retryMaxAttempts = 3
)

// ScheduleParams contains the data needed to generate reminder tasks.
type ScheduleParams struct {
	AppointmentID uuid.UUID
	AppointmentAt time.Time
	CustomerPhone shared.Phone
	EmployeePhone shared.Phone
	ServiceName   string
	LocationName  string
}

// GenerateReminders creates tasks for all applicable reminder rules.
// Skips rules whose scheduled_for is already in the past.
func GenerateReminders(params ScheduleParams) []*Task {
	now := time.Now()

	var tasks []*Task

	metadata := Metadata{
		ServiceName:   params.ServiceName,
		LocationName:  params.LocationName,
		AppointmentAt: params.AppointmentAt,
	}

	for _, rule := range ReminderRules {
		scheduledFor := params.AppointmentAt.Add(-rule.Offset)
		if scheduledFor.Before(now) {
			continue
		}

		tasks = append(tasks, NewTask(CreateTaskParams{
			AppointmentID:  params.AppointmentID,
			Type:           rule.Type,
			RecipientType:  RecipientCustomer,
			RecipientPhone: params.CustomerPhone,
			Metadata:       metadata,
			ScheduledFor:   scheduledFor,
			MaxAttempts:    retryMaxAttempts,
		}))

		tasks = append(tasks, NewTask(CreateTaskParams{
			AppointmentID:  params.AppointmentID,
			Type:           rule.Type,
			RecipientType:  RecipientEmployee,
			RecipientPhone: params.EmployeePhone,
			Metadata:       metadata,
			ScheduledFor:   scheduledFor,
			MaxAttempts:    retryMaxAttempts,
		}))
	}

	return tasks
}
