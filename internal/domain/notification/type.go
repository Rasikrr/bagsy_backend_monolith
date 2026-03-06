package notification

import "time"

type Type string

const (
	Type1HrReminder  Type = "1hr_reminder"
	Type24HrReminder Type = "24hr_reminder"
)

// ReminderRule defines a notification type and its offset before the appointment.
type ReminderRule struct {
	Type   Type
	Offset time.Duration
}

// ReminderRules is the canonical list of all reminder offsets.
// To add a new reminder, append here and add a template in messenger.
var ReminderRules = []ReminderRule{
	{Type: Type24HrReminder, Offset: 24 * time.Hour},
	{Type: Type1HrReminder, Offset: 1 * time.Hour},
}
