package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/dto"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/cockroachdb/errors"
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

func (t *TokenManager) NewAccessToken(user *entity.User, ttl time.Duration) (string, error) {
	claims := t.createAccessClaims(user, ttl)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}
	return tokenStr, nil
}

// ParseAccessToken парсит access токен и возвращает DTO
// Конвертация в domain.Session должна происходить на уровне Service
// nolint: nestif
func (t *TokenManager) ParseAccessToken(accessToken string) (*dto.AccessTokenPayload, error) {
	claims := new(accessClaims)
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrUnexpectedSigningMethod, token.Header["alg"])
		}
		return []byte(t.secretKey), nil
	})
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			}
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrMalformedToken
			}
			if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return nil, ErrInvalidSignature
			}
		}
		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}
	if !token.Valid {
		return nil, ErrTokenNotValid
	}

	// Возвращаем DTO из domain
	return &dto.AccessTokenPayload{
		Phone:       claims.Phone,
		Role:        claims.Role,
		PointCode:   claims.PointCode,
		NetworkCode: claims.NetworkCode,
	}, nil
}

func (t *TokenManager) createAccessClaims(user *entity.User, ttl time.Duration) *accessClaims {
	jwtID := uuid.New().String()
	claims := &accessClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        jwtID,
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    t.issuer,
		},
		Phone: user.Phone,
		Role:  user.Role.String(),
	}
	if user.PointCode != nil {
		claims.PointCode = *user.PointCode
	}
	if user.NetworkCode != nil {
		claims.NetworkCode = *user.NetworkCode
	}
	return claims
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

func (t *TokenManager) NewAuthToken(dto *dto.RegistrationTokenPayload, ttl time.Duration) (string, error) {
	jwtID := uuid.New().String()
	claims := &registrationClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        jwtID,
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    t.issuer,
			Subject:   dto.Phone,
		},
		Phone:       dto.Phone,
		PointCode:   dto.PointCode,
		NetworkCode: dto.NetworkCode,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign registration token: %w", err)
	}
	return tokenStr, nil
}

// nolint: nestif
// ParseAuthToken парсит registration токен и возвращает DTO
func (t *TokenManager) ParseAuthToken(tokenString string) (*dto.RegistrationTokenPayload, error) {
	claims := new(registrationClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrUnexpectedSigningMethod, token.Header["alg"])
		}
		return []byte(t.secretKey), nil
	})
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrRegistrationTokenExpired
			}
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrMalformedToken
			}
			if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return nil, ErrInvalidSignature
			}
		}
		return nil, fmt.Errorf("%w: %w", ErrInvalidRegistrationToken, err)
	}
	if !token.Valid {
		return nil, ErrInvalidRegistrationToken
	}

	// Проверяем наличие обязательных полей
	if claims.Phone == "" {
		return nil, fmt.Errorf("%w: phone", ErrMissingRequiredClaim)
	}
	if claims.PointCode == "" {
		return nil, fmt.Errorf("%w: point_code", ErrMissingRequiredClaim)
	}
	if claims.NetworkCode == "" {
		return nil, fmt.Errorf("%w: network_code", ErrMissingRequiredClaim)
	}

	return &dto.RegistrationTokenPayload{
		Phone:       claims.Phone,
		PointCode:   claims.PointCode,
		NetworkCode: claims.NetworkCode,
	}, nil
}
