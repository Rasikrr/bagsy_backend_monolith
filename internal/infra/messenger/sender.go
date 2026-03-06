package messenger

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/notification"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
)

type MessageSender interface {
	SendMessage(ctx context.Context, phone, message string) error
}

type Messenger struct {
	primary  MessageSender
	fallback MessageSender
}

func NewMessenger(primary, fallback MessageSender) *Messenger {
	return &Messenger{
		primary:  primary,
		fallback: fallback,
	}
}

func (s *Messenger) SendOTP(ctx context.Context, phone shared.Phone, code string) error {
	msg := formatOTPMessage(code)
	return s.send(ctx, phone, msg)
}

func (s *Messenger) SendBookingConfirmationCode(ctx context.Context, phone shared.Phone, code string) error {
	msg := formatBookingOTPMessage(code)
	return s.send(ctx, phone, msg)
}

func (s *Messenger) SendPasswordResetLink(ctx context.Context, phone shared.Phone, link string) error {
	msg := formatPasswordResetMessage(link)
	return s.send(ctx, phone, msg)
}

func (s *Messenger) SendInviteLink(ctx context.Context, phone shared.Phone, link string) error {
	msg := formatInviteMessage(link)
	return s.send(ctx, phone, msg)
}

// SendReminder formats a reminder message based on metadata and sends it.
func (s *Messenger) SendReminder(ctx context.Context, task *notification.Task) error {
	formattedTime := task.Metadata.AppointmentAt.Format("02.01.2006 15:04")
	msg := formatReminderMessage(
		task.Type,
		task.RecipientType,
		task.Metadata.ServiceName,
		task.Metadata.LocationName,
		formattedTime,
	)
	return s.send(ctx, task.RecipientPhone, msg)
}

// FormatReminderMessage formats a reminder message based on type and recipient.
// Implements notification.MessageFormatter.
func formatReminderMessage(taskType notification.Type, recipientType notification.RecipientType, serviceName, locationName, startAt string) string {
	switch {
	case recipientType == notification.RecipientCustomer && taskType == notification.Type24HrReminder:
		return formatCustomer24hReminder(startAt, serviceName, locationName)
	case recipientType == notification.RecipientCustomer && taskType == notification.Type1HrReminder:
		return formatCustomer1hReminder(serviceName, locationName)
	case recipientType == notification.RecipientEmployee && taskType == notification.Type24HrReminder:
		return formatEmployee24hReminder(startAt, serviceName, locationName)
	case recipientType == notification.RecipientEmployee && taskType == notification.Type1HrReminder:
		return formatEmployee1hReminder(serviceName, locationName)
	default:
		return fmt.Sprintf("Напоминание о записи на %s (%s). Bagsy", serviceName, locationName)
	}
}

func (s *Messenger) send(ctx context.Context, phone shared.Phone, msg string) error {
	if err := s.primary.SendMessage(ctx, phone.String(), msg); err != nil {
		return s.fallback.SendMessage(ctx, phone.String(), msg)
	}
	return nil
}
