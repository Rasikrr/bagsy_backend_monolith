package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type TokenManager struct {
	secretKey string
	issuer    string
}

func NewTokenManager(secretKey, issuer string) *TokenManager {
	return &TokenManager{
		secretKey: secretKey,
		issuer:    issuer,
	}
}

func (t *TokenManager) NewAccessToken(authToken auth.Token) (string, error) {
	claims := t.createAccessClaims(authToken)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := jwtToken.SignedString([]byte(t.secretKey))

	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}
	return tokenStr, nil
}

func (t *TokenManager) NewRefreshToken() (raw, hash string, err error) {
	b := make([]byte, 32)
	if _, readErr := rand.Read(b); readErr != nil {
		return "", "", readErr
	}

	raw = base64.RawURLEncoding.EncodeToString(b)

	h := sha256.Sum256([]byte(raw))
	hash = hex.EncodeToString(h[:])

	return raw, hash, nil
}

// ParseAccessToken парсит access токен и возвращает auth.Token
func (t *TokenManager) ParseAccessToken(accessToken string) (auth.Token, error) {
	claims := new(accessClaims)
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.secretKey), nil
	})
	if err != nil {
		return auth.Token{}, fmt.Errorf("failed to parse access token: %w", err)
	}
	if !token.Valid {
		return auth.Token{}, errors.New("invalid token")
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return auth.Token{}, err
	}
	phone, err := shared.NewPhone(claims.Phone)
	if err != nil {
		return auth.Token{}, err
	}
	expiresAt := time.Unix(claims.ExpiresAt, 0)

	return auth.ReconstructToken(userID, phone, expiresAt), nil
}

func (t *TokenManager) createAccessClaims(token auth.Token) *accessClaims {
	jwtID := uuid.New().String()
	claims := &accessClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        jwtID,
			ExpiresAt: token.ExpiresAt.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    t.issuer,
		},
		Phone:  token.Phone.String(),
		UserID: token.UserID.String(),
	}
	return claims
}
