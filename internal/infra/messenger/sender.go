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
	if err := s.primary.SendMessage(ctx, phone.String(), msg); err != nil {
		return s.fallback.SendMessage(ctx, phone.String(), msg)
	}
	return nil
}
