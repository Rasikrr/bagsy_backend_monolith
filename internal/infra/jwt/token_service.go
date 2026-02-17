package jwt

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

type tokenManager interface {
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
	tokenGenerator   tokenManager
	refreshTokenRepo refreshTokenRepository
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
}

func NewTokenService(
	tokenGenerator tokenManager,
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

func (s *TokenService) GenerateTokens(ctx context.Context, userID uuid.UUID, phone shared.Phone) (access, refresh string, err error) {
	tokenInfo := auth.NewToken(userID, phone, s.accessTokenTTL)
	accessToken, err := s.tokenGenerator.NewAccessToken(tokenInfo)
	if err != nil {
		//TODO: handler error (business error)
		return "", "", err
	}
	refreshToken, refreshHash, err := s.tokenGenerator.NewRefreshToken()
	if err != nil {
		//TODO: handler error (business error)
		return "", "", err
	}
	err = s.refreshTokenRepo.SaveToken(ctx, refreshHash, userID, s.refreshTokenTTL)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *TokenService) VerifyAccessToken(_ context.Context, tokenStr string) (*auth.Token, error) {
	tokenInfo, err := s.tokenGenerator.ParseAccessToken(tokenStr)
	if err != nil {
		// TODO: error
		return nil, err
	}
	return &tokenInfo, nil
}

func (s *TokenService) DeleteRefreshToken(ctx context.Context, token string) error {
	h := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(h[:])

	return s.refreshTokenRepo.DeleteToken(ctx, tokenHash)
}
