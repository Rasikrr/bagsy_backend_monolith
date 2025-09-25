package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/Rasikrr/core/database"
	"github.com/Rasikrr/core/log"
	"github.com/jackc/pgx/v5"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/util/hash"
	"github.com/Rasikrr/core/version"

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
	RegisterConfirm(ctx context.Context, phone string, password string) (*entity.Auth, error)
	RefreshTokens(ctx context.Context, token string) (*entity.Auth, error)
	ValidateToken(ctx context.Context, token string) (bool, error)
	GetAuthTokenPayload(ctx context.Context, token string) (*entity.PayloadParams, error)
	Login(ctx context.Context, phone string, password string) (*entity.Auth, error)
	GenAuthConfirmationLink(ctx context.Context, phone, pointCode string) (string, error)
}

type service struct {
	smsClient   sms.Client
	authCache   auth.Cache
	tgClient    telegram.Client
	userService users.Service
	txManager   database.TXManager

	tgChatID            int64
	jwtSecret           string
	authConfirmationURL string
}

func NewService(
	smsClient sms.Client,
	tgClient telegram.Client,
	authCache auth.Cache,
	userService users.Service,
	txManager database.TXManager,
	tgChatID int64,
	authConfirmationURL string,
	jwtSecret string,
) Service {
	return &service{
		smsClient:           smsClient,
		tgClient:            tgClient,
		authCache:           authCache,
		userService:         userService,
		tgChatID:            tgChatID,
		txManager:           txManager,
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

	log.Infof(ctx, "user %+v", user)

	if err != nil {
		return nil, err
	}
	// TODO: errors
	if user.Password == nil {
		return nil, errNoAccess
	}
	if !user.Active {
		return nil, errNoAccess
	}
	valid := hash.CheckPassword(*user.Password, password)
	if !valid {
		return nil, errInvalidPassword
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

func (s *service) GenAuthConfirmationLink(_ context.Context, phone, pointCode string) (string, error) {
	token, err := jwt.GenerateRegistrationToken(phone, pointCode, s.jwtSecret)
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

func (s *service) RegisterConfirm(ctx context.Context, phone string, password string) (*entity.Auth, error) {
	log.Infof(ctx, "phone %v", phone)

	var result *entity.Auth

	err := s.txManager.Transaction(ctx, pgx.TxOptions{}, func(txCtx context.Context) error {
		user, err := s.userService.GetByPhone(txCtx, phone)
		if err != nil {
			return fmt.Errorf("get user: %w", err)
		}

		hashedPassword, err := hash.Password(password)
		if err != nil {
			return fmt.Errorf("hashing failed: %w", err)
		}

		err = s.userService.SetPasswordByPhone(txCtx, user.Phone, hashedPassword)
		if err != nil {
			return fmt.Errorf("set password: %w", err)
		}

		err = s.userService.SetActive(txCtx, phone)
		if err != nil {
			return fmt.Errorf("activate user: %w", err)
		}

		updatedUser, err := s.userService.GetByPhone(txCtx, phone)
		if err != nil {
			return fmt.Errorf("get updated user: %w", err)
		}

		accessToken, refreshToken, err := s.generateTokens(updatedUser)
		if err != nil {
			return fmt.Errorf("generate tokens: %w", err)
		}

		result = &entity.Auth{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("registration confirm failed: %w", err)
	}
	return result, nil
}

func (s *service) ValidateToken(_ context.Context, token string) (bool, error) {
	return jwt.ValidateToken(token, s.jwtSecret)
}

func (s *service) GetAuthTokenPayload(_ context.Context, token string) (*entity.PayloadParams, error) {
	return jwt.ParseAuthToken(token, s.jwtSecret)
}

func (s *service) RefreshTokens(ctx context.Context, token string) (*entity.Auth, error) {
	valid, err := jwt.ValidateToken(token, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("validate token: %w", err)
	}
	if !valid {
		return nil, errors.New("invalid token")
	}
	payload, err := jwt.ParseRefreshToken(token, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("parse refresh token: %w", err)
	}
	if !payload.IsRefresh() {
		return nil, errors.New("access token is not allowed")
	}
	user, err := s.userService.GetByPhone(ctx, payload.Phone)
	if err != nil {
		return nil, fmt.Errorf("get user by phone: %w", err)
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

	refreshToken, err = jwt.GenerateRefreshToken(refreshParams, s.jwtSecret)
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}
