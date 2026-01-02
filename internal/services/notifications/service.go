package notifications

import (
	"context"
	"fmt"
)

type smsClient interface {
	Send(ctx context.Context, phone, message string) error
}

type whatsappClient interface {
	SendMessage(_ context.Context, phoneNumber, message string) error
}

type Service struct {
	smsClient              smsClient
	whatsApp               whatsappClient
	registrationConfirmURL string
}

func NewService(
	smsClient smsClient,
	whatsApp whatsappClient,
	registrationConfirmURL string,
) *Service {
	return &Service{
		smsClient:              smsClient,
		whatsApp:               whatsApp,
		registrationConfirmURL: registrationConfirmURL,
	}
}

func (s *Service) SendRegistrationLink(ctx context.Context, phone, token string) error {
	link := fmt.Sprintf("%s?token=%s", s.registrationConfirmURL, token)
	// TODO: format message (markdown)
	message := fmt.Sprintf("Добро пожаловать в Bagsy! Завершите регистрацию по ссылке: %s", link)

	err := s.whatsApp.SendMessage(ctx, phone, message)
	if err != nil {
		err = s.smsClient.Send(ctx, phone, message)
		if err != nil {
			return err
		}
	}
	return nil
}
