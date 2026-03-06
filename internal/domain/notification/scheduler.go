package notification

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

// MessageFormatter formats a reminder message based on notification type, recipient, and appointment details.
type MessageFormatter func(taskType Type, recipientType RecipientType, serviceName, locationName, startAt string) string

// ScheduleParams contains the data needed to generate reminder tasks.
type ScheduleParams struct {
	AppointmentID uuid.UUID
	AppointmentAt time.Time
	CustomerPhone shared.Phone
	EmployeePhone shared.Phone
	ServiceName   string
	LocationName  string
	MaxAttempts   int
	Formatter     MessageFormatter
}

// GenerateReminders creates tasks for all applicable reminder rules.
// Skips rules whose scheduled_for is already in the past.
func GenerateReminders(params ScheduleParams) []*Task {
	now := time.Now()
	formattedTime := params.AppointmentAt.Format("02.01.2006 15:04")

	var tasks []*Task

	for _, rule := range ReminderRules {
		scheduledFor := params.AppointmentAt.Add(-rule.Offset)
		if scheduledFor.Before(now) {
			continue
		}

		customerMsg := params.Formatter(rule.Type, RecipientCustomer, params.ServiceName, params.LocationName, formattedTime)

		tasks = append(tasks, NewTask(CreateTaskParams{
			AppointmentID:  params.AppointmentID,
			Type:           rule.Type,
			RecipientType:  RecipientCustomer,
			RecipientPhone: params.CustomerPhone,
			Message:        customerMsg,
			ScheduledFor:   scheduledFor,
			MaxAttempts:    params.MaxAttempts,
		}))

		employeeMsg := params.Formatter(rule.Type, RecipientEmployee, params.ServiceName, params.LocationName, formattedTime)

		tasks = append(tasks, NewTask(CreateTaskParams{
			AppointmentID:  params.AppointmentID,
			Type:           rule.Type,
			RecipientType:  RecipientEmployee,
			RecipientPhone: params.EmployeePhone,
			Message:        employeeMsg,
			ScheduledFor:   scheduledFor,
			MaxAttempts:    params.MaxAttempts,
		}))
	}

	return tasks
}
