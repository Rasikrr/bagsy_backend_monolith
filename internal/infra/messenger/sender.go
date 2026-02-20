package messenger

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
)

type MessageSender interface {
	SendMessage(ctx context.Context, phone, message string) error
}

type OTPSender struct {
	primary  MessageSender
	fallback MessageSender
}

func NewOTPSender(primary, fallback MessageSender) *OTPSender {
	return &OTPSender{
		primary:  primary,
		fallback: fallback,
	}
}

func (s *OTPSender) SendOTP(ctx context.Context, phone shared.Phone, code string) error {
	msg := formatOTPMessage(code)
	return s.send(ctx, phone, msg)
}

func (s *OTPSender) SendBookingConfirmationCode(ctx context.Context, phone shared.Phone, code string) error {
	msg := formatBookingOTPMessage(code)
	return s.send(ctx, phone, msg)
}

func (s *OTPSender) SendPasswordResetLink(ctx context.Context, phone shared.Phone, link string) error {
	msg := formatPasswordResetMessage(link)
	return s.send(ctx, phone, msg)
}

func (s *OTPSender) SendInviteLink(ctx context.Context, phone shared.Phone, link string) error {
	msg := formatInviteMessage(link)
	return s.send(ctx, phone, msg)
}

func (s *OTPSender) send(ctx context.Context, phone shared.Phone, msg string) error {
	if err := s.primary.SendMessage(ctx, phone.String(), msg); err != nil {
		return s.fallback.SendMessage(ctx, phone.String(), msg)
	}
	return nil
}
