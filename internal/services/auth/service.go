package auth

import (
	"context"
	"fmt"

	"github.com/Rasikrr/core/version"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/cache/auth"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/clients/sms"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/services/users"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/util/codegen"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/util/hash"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/util/jwt"
	"github.com/Rasikrr/core/enum"
	"github.com/Rasikrr/core/telegram"
)

type Service interface {
	SendCode(ctx context.Context, phone string) error
	Login(ctx context.Context, phone string, password string) (*entity.Auth, error)
	GenAuthConfirmationLink(ctx context.Context, phone string) (string, error)
}

type service struct {
	smsClient   sms.Client
	authCache   auth.Cache
	tgClient    telegram.Client
	userService users.Service

	tgChatID            int64
	jwtSecret           string
	authConfirmationURL string
}

func NewService(
	smsClient sms.Client,
	tgClient telegram.Client,
	authCache auth.Cache,
	userService users.Service,
	tgChatID int64,
	authConfirmationURL string,
	jwtSecret string,
) Service {
	return &service{
		smsClient:   smsClient,
		tgClient:    tgClient,
		authCache:   authCache,
		userService: userService,
		tgChatID:    tgChatID,

		authConfirmationURL: authConfirmationURL,
		jwtSecret:           jwtSecret,
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

	if version.GetVersion() != enum.EnvironmentProd {
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

	hashedPassword, err := hash.Password(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	err = s.userService.SetPasswordByPhone(ctx, phone, hashedPassword)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return &entity.Auth{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) GenAuthConfirmationLink(_ context.Context, phone string) (string, error) {
	token, err := jwt.GenerateRegistrationToken(phone, s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("generate registration token: %w", err)
	}
	return fmt.Sprintf("%s?token=%s", s.authConfirmationURL, token), nil
}

func (s *service) prepareMessage(phone, code string) string {
	if version.GetVersion() != enum.EnvironmentProd {
		return fmt.Sprintf("%s: %s - код для входа на bagsy.kz", phone, code)
	}
	return fmt.Sprintf("%s: Ваш код для входа в bagsy.kz", code)
}

// nolint: nonamedreturns
func (s *service) generateTokens(user *entity.User) (accessToken, refreshToken string, err error) {
	accessParams := &entity.PayloadParams{
		Phone:   user.Phone,
		Role:    user.Role.String(),
		Active:  user.Active,
		Refresh: false,
	}

	accessToken, err = jwt.GenerateAccessToken(accessParams, s.jwtSecret)
	if err != nil {
		return "", "", fmt.Errorf("generate access token: %w", err)
	}

	refreshParams := &entity.PayloadParams{
		Phone:   user.Phone,
		Role:    user.Role.String(),
		Active:  user.Active,
		Refresh: true,
	}

	refreshToken, err = jwt.GenerateRefreshToken(refreshParams)
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}
