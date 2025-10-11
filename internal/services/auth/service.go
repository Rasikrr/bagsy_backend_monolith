package auth

import (
	"context"
	"fmt"

	"github.com/Rasikrr/core/database"
	"github.com/Rasikrr/core/log"
	"github.com/jackc/pgx/v5"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/hash"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/cache/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/sms"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/whatsapp"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/jwt"
	"github.com/Rasikrr/core/telegram"
)

type Service interface {
	SendRegisterLink(ctx context.Context, phone, link string) error
	RegisterConfirm(ctx context.Context, phone string, password string) (*entity.Auth, error)
	RefreshTokens(ctx context.Context, token string) (*entity.Auth, error)
	ValidateRegisterToken(ctx context.Context, token string) (bool, error)
	GetAuthTokenPayload(ctx context.Context, token string) (*entity.PayloadParams, error)
	Login(ctx context.Context, phone string, password string) (*entity.Auth, error)
	GenAuthConfirmationLink(ctx context.Context, phone, pointCode string) (string, error)
}

type service struct {
	smsClient      sms.Client
	whatsAppClient whatsapp.Client
	authCache      auth.Cache
	tgClient       telegram.Client
	userService    users.Service
	txManager      database.TXManager

	tgChatID            int64
	jwtSecret           string
	authConfirmationURL string
}

func NewService(
	smsClient sms.Client,
	whatsAppClient whatsapp.Client,
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
		whatsAppClient:      whatsAppClient,
		tgClient:            tgClient,
		authCache:           authCache,
		userService:         userService,
		tgChatID:            tgChatID,
		txManager:           txManager,
		authConfirmationURL: authConfirmationURL,
		jwtSecret:           jwtSecret,
	}
}

func (s *service) SendRegisterLink(ctx context.Context, phone, link string) error {
	message := fmt.Sprintf("Ваша ссылка на регистрацию на bagsy.kz: %s", link)
	return s.whatsAppClient.SendMessage(ctx, phone, message)
}

func (s *service) Login(ctx context.Context, phone string, password string) (*entity.Auth, error) {
	user, err := s.userService.GetByPhone(ctx, phone)

	if err != nil {
		return nil, err
	}
	// TODO: errorsw
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
		return nil, errGenerateTokens.Wrap(err)
	}

	return &entity.Auth{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) GenAuthConfirmationLink(_ context.Context, phone, pointCode string) (string, error) {
	token, err := jwt.GenerateRegistrationToken(phone, pointCode, s.jwtSecret)
	if err != nil {
		return "", errGenerateRegistrationURL.Wrap(err)
	}
	return fmt.Sprintf("%s?token=%s", s.authConfirmationURL, token), nil
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
			return errHashingFailed.Wrap(err)
		}

		err = s.userService.SetPasswordByPhone(txCtx, user.Phone, hashedPassword)
		if err != nil {
			return errSetPassword.Wrap(err)
		}

		err = s.userService.SetActive(txCtx, phone)
		if err != nil {
			return errActivateUser.Wrap(err)
		}

		updatedUser, err := s.userService.GetByPhone(txCtx, phone)
		if err != nil {
			return fmt.Errorf("get updated user: %w", err)
		}

		accessToken, refreshToken, err := s.generateTokens(updatedUser)
		if err != nil {
			return errGenerateTokens.Wrap(err)
		}

		result = &entity.Auth{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		return nil
	})

	if err != nil {
		return nil, errRegistrationConfirm.Wrap(err)
	}
	return result, nil
}

func (s *service) ValidateRegisterToken(_ context.Context, token string) (bool, error) {
	return jwt.ValidateToken(token, s.jwtSecret)
}

func (s *service) GetAuthTokenPayload(_ context.Context, token string) (*entity.PayloadParams, error) {
	return jwt.ParseAuthToken(token, s.jwtSecret)
}

func (s *service) RefreshTokens(ctx context.Context, token string) (*entity.Auth, error) {
	valid, err := jwt.ValidateToken(token, s.jwtSecret)
	if err != nil {
		return nil, errInvalidToken.Wrap(err)
	}
	if !valid {
		return nil, errInvalidToken
	}
	payload, err := jwt.ParseRefreshToken(token, s.jwtSecret)
	if err != nil {
		return nil, errInvalidToken.Wrap(err)
	}
	if !payload.IsRefresh() {
		return nil, errAccessTokenNotAllowed
	}
	user, err := s.userService.GetByPhone(ctx, payload.Phone)
	if err != nil {
		return nil, err
	}
	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, errGenerateTokens.Wrap(err)
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
	if user.PointCode != nil {
		accessParams.PointCode = *user.PointCode
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
	if user.PointCode != nil {
		refreshParams.PointCode = *user.PointCode
	}

	refreshToken, err = jwt.GenerateRefreshToken(refreshParams, s.jwtSecret)
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}
