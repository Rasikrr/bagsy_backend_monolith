package notifications

import (
	"context"
	"fmt"

	"github.com/Rasikrr/core/log"
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

func (s *Service) SendManagementAuthConfirmationCode(ctx context.Context, phone, code string) error {
	// TODO: format message (markdown)
	message := fmt.Sprintf("%s - Ваш код подтверждения регистрации в системе Bagsy", code)
	return s.send(ctx, phone, message)
}

func (s *Service) SendStaffRegistrationLink(ctx context.Context, phone, token string) error {
	link := fmt.Sprintf("%s/%s", s.registrationConfirmURL, token)
	// TODO: format message (markdown)
	message := fmt.Sprintf("Добро пожаловать в Bagsy! Завершите регистрацию по ссылке: %s", link)
	return s.send(ctx, phone, message)
}

func (s *Service) SendBagsyConfirmCode(ctx context.Context, phone, code string) error {
	// TODO: format message, add link, name of service etc. (markdown)
	message := fmt.Sprintf("%s - Ваш код подтверждения на запись", code)
	return s.send(ctx, phone, message)
}

func (s *Service) SendPasswordChangeLink(ctx context.Context, phone, token string) error {
	link := fmt.Sprintf("%s/%s", s.registrationConfirmURL, token)
	// TODO: format message (markdown)
	message := fmt.Sprintf("Для смены пароля в Bagsy следуйте по данной ссылке: %s", link)
	return s.send(ctx, phone, message)
}

func (s *Service) send(ctx context.Context, phone, message string) error {
	err := s.whatsApp.SendMessage(ctx, phone, message)
	if err != nil {
		log.Warnf(ctx, "Failed to send message by whatsapp: %v", err)
		err = s.smsClient.Send(ctx, phone, message)
		if err != nil {
			return err
		}
	}
	return nil
}
