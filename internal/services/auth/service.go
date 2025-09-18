package auth

import (
	"context"
	"fmt"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/cache/auth"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/clients/sms"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/services/users"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/util/codegen"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/util/jwt"
	"github.com/Rasikrr/core/enum"
	"github.com/Rasikrr/core/telegram"
)

type Service interface {
	SendCode(ctx context.Context, phone string) error
	Login(ctx context.Context, phone string, password string) (*entity.Auth, error)
}

type service struct {
	env enum.Environment

	smsClient   sms.Client
	authCache   auth.Cache
	tgClient    telegram.Client
	userService users.Service

	tgChatID int64
}

func NewService(
	env enum.Environment,
	smsClient sms.Client,
	tgClient telegram.Client,
	authCache auth.Cache,
	userService users.Service,
	tgChatID int64,
) Service {
	return &service{
		smsClient:   smsClient,
		tgClient:    tgClient,
		authCache:   authCache,
		userService: userService,
		tgChatID:    tgChatID,

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

func (s *service) Login(ctx context.Context, phone string, password string) (*entity.Auth, error) {
	user, err := s.userService.GetByPhone(ctx, phone)
	if err != nil {
		return nil, errUserNotFound
	}

	hashedPassword, err := jwt.HashPassword(password)

	err = s.userService.SetPasswordByPhone(ctx, phone, hashedPassword)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken := s.generateTokens(user)

	return &entity.Auth{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) prepareMessage(phone, code string) string {
	if s.env == enum.EnvironmentDev {
		return fmt.Sprintf("%s: %s - код для входа на bagsy.kz", phone, code)
	}
	return fmt.Sprintf("%s: Ваш код для входа в bagsy.kz", code)
}

func (s *service) generateTokens(user *entity.User) (accessToken string, refreshToken string) {
	accessParams := &entity.PayloadParams{
		Phone:   user.Phone,
		Role:    user.Role.String(),
		Active:  user.Active,
		Refresh: false,
	}

	accessToken, err := jwt.GenerateAccessToken(accessParams)
	if err != nil {
		return "", ""
	}

	refreshParams := &entity.PayloadParams{
		Phone:   user.Phone,
		Role:    user.Role.String(),
		Active:  user.Active,
		Refresh: true,
	}

	refreshToken, err = jwt.GenerateRefreshToken(refreshParams)
	if err != nil {
		return "", ""
	}

	return accessToken, refreshToken
}
