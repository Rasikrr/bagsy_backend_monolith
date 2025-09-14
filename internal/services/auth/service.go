package auth

import (
	"context"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/clients/sms"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/core/telegram"
)

type Service interface {
	Login(ctx context.Context, phone string) (*entity.Auth, error)
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

func (s *service) Login(ctx context.Context, phone string) (*entity.Auth, error) {
	// TODO
	return nil, nil
}
