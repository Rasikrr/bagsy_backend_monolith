package jwt

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
	GetToken(ctx context.Context, tokenHash string) (userID uuid.UUID, err error)
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
		return "", "", fmt.Errorf("generate access token: %w", err)
	}
	refreshToken, refreshHash, err := s.tokenGenerator.NewRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}
	if err := s.refreshTokenRepo.SaveToken(ctx, refreshHash, userID, s.refreshTokenTTL); err != nil {
		return "", "", fmt.Errorf("save refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *TokenService) VerifyAccessToken(_ context.Context, tokenStr string) (*auth.Token, error) {
	tokenInfo, err := s.tokenGenerator.ParseAccessToken(tokenStr)
	if err != nil {
		return nil, fmt.Errorf("verify access token: %w", err)
	}
	return &tokenInfo, nil
}

func (s *TokenService) RefreshTokens(ctx context.Context, oldRefreshToken string) (userID uuid.UUID, err error) {
	h := sha256.Sum256([]byte(oldRefreshToken))
	oldHash := hex.EncodeToString(h[:])

	// 1. Validate old refresh token, get userID.
	userID, err = s.refreshTokenRepo.GetToken(ctx, oldHash)
	if err != nil {
		return uuid.Nil, fmt.Errorf("get refresh token: %w", err)
	}

	// 2. Delete old refresh token (rotation).
	if err := s.refreshTokenRepo.DeleteToken(ctx, oldHash); err != nil {
		return uuid.Nil, fmt.Errorf("delete old refresh token: %w", err)
	}

	return userID, nil
}

func (s *TokenService) DeleteRefreshToken(ctx context.Context, token string) error {
	h := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(h[:])

	return s.refreshTokenRepo.DeleteToken(ctx, tokenHash)
}
