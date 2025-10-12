package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
)

const (
	accessTTL  = 15 * time.Minute
	refreshTTL = 3 * 24 * time.Hour
)

var (
	errInvalidClaims     = errors.New("unable to parse claims")
	errTokenNotValid     = errors.New("token is not valid")
	errUnexpectedSigning = errors.New("unexpected signing method")
)

func GenerateAccessToken(params *entity.PayloadParams, secret string) (string, error) {
	params.ExpiresAt = time.Now().Add(accessTTL).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, params)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(params *entity.PayloadParams, secret string) (string, error) {
	params.ExpiresAt = time.Now().Add(refreshTTL).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, params)
	return token.SignedString([]byte(secret))
}

func ValidateToken(token string, secret string) (bool, error) {
	parsedToken, err := parseToken(token, secret)
	if err != nil {
		return false, err
	}
	return parsedToken.Valid, nil
}

func ParseAuthToken(tokenString string, secret string) (*entity.PayloadParams, error) {
	return parseTokenWithClaims(tokenString, secret, true)
}

func ParseRefreshToken(tokenString string, secret string) (*entity.PayloadParams, error) {
	return parseTokenWithClaims(tokenString, secret, false)
}

// parseToken - общая функция парсинга токена
func parseToken(tokenString string, secret string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", errUnexpectedSigning, token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errTokenNotValid
	}

	return parsedToken, nil
}

// parseTokenWithClaims - парсит токен и извлекает claims
func parseTokenWithClaims(tokenString string, secret string, includePointData bool) (*entity.PayloadParams, error) {
	parsedToken, err := parseToken(tokenString, secret)
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errInvalidClaims
	}

	return extractClaims(claims, includePointData)
}

// extractClaims - извлекает данные из claims
func extractClaims(claims jwt.MapClaims, includePointData bool) (*entity.PayloadParams, error) {
	phone, ok := claims["phone"].(string)
	if !ok {
		return nil, fmt.Errorf("%w: phone field missing or invalid", errInvalidToken)
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, fmt.Errorf("%w: role field missing or invalid", errInvalidToken)
	}

	active, ok := claims["active"].(bool)
	if !ok {
		return nil, fmt.Errorf("%w: active field missing or invalid", errInvalidToken)
	}

	refresh, ok := claims["refresh"].(bool)
	if !ok {
		return nil, fmt.Errorf("%w: refresh field missing or invalid", errInvalidToken)
	}

	params := &entity.PayloadParams{
		Phone:   phone,
		Role:    role,
		Active:  active,
		Refresh: refresh,
	}
	// nolint: govet
	// Для auth токена дополнительно извлекаем point и network
	if includePointData {
		pointCode, ok := claims["point_code"].(string)
		if !ok {
			return nil, fmt.Errorf("%w: point_code field missing or invalid", errInvalidToken)
		}

		networkCode, ok := claims["network_code"].(string)
		if !ok {
			return nil, fmt.Errorf("%w: network_code field missing or invalid", errInvalidToken)
		}

		params.PointCode = pointCode
		params.NetworkCode = networkCode
	}

	return params, nil
}
