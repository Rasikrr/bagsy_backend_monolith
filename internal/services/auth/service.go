package auth

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/cache/auth"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/clients/sms"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/util/codegen"
	"github.com/Rasikrr/core/enum"
	"github.com/Rasikrr/core/telegram"
)

type Service interface {
	SendCode(ctx context.Context, phone string) error
}

type service struct {
	env enum.Environment

	smsClient sms.Client
	authCache auth.Cache
	tgClient  telegram.Client

	tgChatID int64
}

func NewService(
	env enum.Environment,
	smsClient sms.Client,
	tgClient telegram.Client,
	authCache auth.Cache,
	tgChatID int64,
) Service {
	return &service{
		smsClient: smsClient,
		tgClient:  tgClient,
		authCache: authCache,
		tgChatID:  tgChatID,

		env: env,
	}
}

func (s *service) SendCode(ctx context.Context, phone string) (err error) {
	spam, err := s.authCache.CheckSpam(ctx, phone)
	if err != nil {
		return fmt.Errorf("check spam: %w", err)
	}
	if spam {
		return errSpam
	}
	code := codegen.GenerateAuthCode()
	msg := s.prepareMessage(phone, code)

	defer func() {
		err = s.authCache.SetCode(ctx, phone, code)
	}()

	if s.env == enum.EnvironmentDev {
		return s.tgClient.SendText(ctx, s.tgChatID, msg)
	}
	err = s.smsClient.Send(ctx, phone, msg)
	if err != nil {
		return fmt.Errorf("send sms: %w", err)
	}
	return nil
}

func (s *service) prepareMessage(phone, code string) string {
	if s.env == enum.EnvironmentDev {
		return fmt.Sprintf("%s: %s - код для входа на bagsy.kz", phone, code)
	}
	return fmt.Sprintf("%s: Ваш код для входа в bagsy.kz", code)
}
