package messenger

import (
	"context"

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

func (s *Messenger) send(ctx context.Context, phone shared.Phone, msg string) error {
	if err := s.primary.SendMessage(ctx, phone.String(), msg); err != nil {
		return s.fallback.SendMessage(ctx, phone.String(), msg)
	}
	return nil
}
