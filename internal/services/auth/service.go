package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/Rasikrr/core/database"
	"github.com/jackc/pgx/v5"

	"github.com/Rasikrr/bagsy_backend_monolith/pkg/hash"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/cache/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/sms"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/whatsapp"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/jwt"
	"github.com/Rasikrr/core/telegram"
)

type Service interface {
	Login(ctx context.Context, phone string, password string) (string, string, error)
	CheckAccessToken(ctx context.Context, token string) (*entity.Session, error)
	RefreshTokens(ctx context.Context, token string) (string, string, error)

	GenRegisterConfirmLink(ctx context.Context, phone, pointCode, networkCode string) (string, error)
	SendRegisterConfirmLink(ctx context.Context, phone, link string) error

	RegisterConfirm(ctx context.Context, phone string, password string) (string, string, error)
	ValidateRegisterToken(ctx context.Context, token string) error
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
	accessTokenTTL      time.Duration
	refreshTokenTTL     time.Duration
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
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
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
		accessTokenTTL:      accessTokenTTL,
		refreshTokenTTL:     refreshTokenTTL,
	}
}

func (s *service) SendRegisterConfirmLink(ctx context.Context, phone, link string) error {
	message := fmt.Sprintf("Ваша ссылка на регистрацию на bagsy.kz: %s", link)
	return s.whatsAppClient.SendMessage(ctx, phone, message)
}

func (s *service) Login(ctx context.Context, phone string, password string) (string, string, error) {
	user, err := s.userService.GetByPhone(ctx, phone)
	if err != nil {
		return "", "", errGetUser.Wrap(err)
	}
	if user.Password == nil {
		return "", "", errNoAccess
	}
	if !user.Active {
		return "", "", errNoAccess
	}
	valid := hash.CheckPassword(*user.Password, password)
	if !valid {
		return "", "", errInvalidPassword
	}
	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return "", "", errGenerateTokens.Wrap(err)
	}
	return accessToken, refreshToken, nil
}

func (s *service) GenRegisterConfirmLink(_ context.Context, phone, pointCode, networkCode string) (string, error) {
	token, err := jwt.GenerateRegistrationToken(phone, pointCode, networkCode, s.jwtSecret)
	if err != nil {
		return "", errGenerateRegistrationURL.Wrap(err)
	}
	return fmt.Sprintf("%s?token=%s", s.authConfirmationURL, token), nil
}

func (s *service) RegisterConfirm(ctx context.Context, phone string, password string) (string, string, error) {
	var (
		access, refresh string
	)
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

		access, refresh, err = s.generateTokens(updatedUser)
		if err != nil {
			return errGenerateTokens.Wrap(err)
		}
		return nil
	})

	if err != nil {
		return "", "", errRegistrationConfirm.Wrap(err)
	}
	return access, refresh, nil
}

func (s *service) ValidateRegisterToken(_ context.Context, token string) error {
	return jwt.ValidateRegistrationToken(token, s.jwtSecret)
}

func (s *service) CheckAccessToken(_ context.Context, token string) (*entity.Session, error) {
	claims, err := jwt.ParseToken(token, s.jwtSecret)
	if err != nil {
		return nil, errInvalidToken.Wrap(err)
	}
	if claims.Refresh {
		return nil, errRefreshTokenNotAllowed
	}
	return claims.ToSession()
}

func (s *service) RefreshTokens(ctx context.Context, token string) (string, string, error) {
	claims, err := jwt.ParseToken(token, s.jwtSecret)
	if err != nil {
		return "", "", errInvalidToken.Wrap(err)
	}
	if !claims.Refresh {
		return "", "", errAccessTokenNotAllowed
	}

	user, err := s.userService.GetByPhone(ctx, claims.Phone)
	if err != nil {
		return "", "", err
	}
	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return "", "", errGenerateTokens.Wrap(err)
	}
	return accessToken, refreshToken, nil
}

// nolint: nonamedreturns
func (s *service) generateTokens(user *entity.User) (accessToken, refreshToken string, err error) {
	accessClaims := jwt.NewClaims(user, s.accessTokenTTL, false)
	accessToken, err = jwt.GenerateToken(accessClaims, s.jwtSecret)
	if err != nil {
		return "", "", errGenerateTokens.Wrap(err)
	}
	refreshClaims := jwt.NewClaims(user, s.refreshTokenTTL, true)
	refreshToken, err = jwt.GenerateToken(refreshClaims, s.jwtSecret)
	if err != nil {
		return "", "", errGenerateTokens.Wrap(err)
	}
	return accessToken, refreshToken, nil
}
