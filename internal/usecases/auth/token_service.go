package auth

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

type tokenGenerator interface {
	NewAccessToken(authToken auth.Token) (string, error)
	NewRefreshToken() (raw, hash string, err error)
	ParseAccessToken(accessToken string) (auth.Token, error)
}

type refreshTokenRepository interface {
	SaveToken(ctx context.Context, tokenHash string, userID uuid.UUID, ttl time.Duration) error
	GetToken(ctx context.Context, tokenHash string) (string, error)
	DeleteToken(ctx context.Context, tokenHash string) error
}

//type tokensRepository interface {
//	SaveInviteToken(ctx context.Context, token string, payload *InviteTokenInfo, ttl time.Duration) error
//	GetInviteToken(ctx context.Context, token string) (*InviteTokenInfo, error)
//	DeleteInviteToken(ctx context.Context, token string) error
//}

type TokenService struct {
	tokenGenerator   tokenGenerator
	refreshTokenRepo refreshTokenRepository
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
}

func NewTokenService(
	tokenGenerator tokenGenerator,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
	refreshTokenRepo refreshTokenRepository,
) *TokenService {
	return &TokenService{
		tokenGenerator:   tokenGenerator,
		accessTokenTTL:   accessTokenTTL,
		refreshTokenTTL:  refreshTokenTTL,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (t *TokenService) GenerateTokens(ctx context.Context, userID uuid.UUID, phone shared.Phone) (access, refresh string, err error) {
	tokenInfo := auth.NewToken(userID, phone, t.accessTokenTTL)
	accessToken, err := t.tokenGenerator.NewAccessToken(tokenInfo)
	if err != nil {
		//TODO: handler error (business error)
		return "", "", err
	}
	refreshToken, refreshHash, err := t.tokenGenerator.NewRefreshToken()
	if err != nil {
		//TODO: handler error (business error)
		return "", "", err
	}
	err = t.refreshTokenRepo.SaveToken(ctx, refreshHash, userID, t.refreshTokenTTL)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
