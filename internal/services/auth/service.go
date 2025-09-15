package auth

import (
	"context"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/clients/sms"
	"github.com/Rasikrr/core/telegram"
)

type Service interface {
	SendCode(ctx context.Context, phone string) error
}

type service struct {
	smsClient sms.Client
	tgClient  telegram.Client
}

func NewService(
	tgClient telegram.Client,
	smsClient sms.Client,
) Service {
	return &service{
		tgClient:  tgClient,
		smsClient: smsClient,
	}
}

func (s *service) SendCode(ctx context.Context, phone string) error {
	// TODO
	return nil
}
