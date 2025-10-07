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

func GenerateAccessToken(params *entity.PayloadParams, secret string) (string, error) {
	claims := jwt.MapClaims{
		"phone":   params.Phone,
		"role":    params.Role,
		"active":  params.Active,
		"refresh": params.Refresh,
		"exp":     time.Now().Add(accessTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(params *entity.PayloadParams, secret string) (string, error) {
	claims := jwt.MapClaims{
		"phone":      params.Phone,
		"refresh":    params.Refresh,
		"point_code": params.PointCode,
		"exp":        time.Now().Add(refreshTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(token string, secret string) (bool, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return false, err
	}

	if !parsedToken.Valid {
		return false, errors.New("token is not valid")
	}

	return true, nil
}

// nolint: govet
func ParseAuthToken(tokenString string, secret string) (*entity.PayloadParams, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("token is not valid")
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		phone, ok := claims["phone"].(string)
		if !ok {
			return nil, errInvalidToken
		}
		role, ok := claims["role"].(string)
		if !ok {
			return nil, errInvalidToken
		}
		refresh, ok := claims["refresh"].(bool)
		if !ok {
			return nil, errInvalidToken
		}
		pointCode, ok := claims["point_code"].(string)
		if !ok {
			return nil, errInvalidToken
		}
		return &entity.PayloadParams{
			Phone:     phone,
			Role:      role,
			Refresh:   refresh,
			PointCode: pointCode,
		}, nil
	}

	return nil, errors.New("unable to parse claims")
}

// nolint: govet
func ParseRefreshToken(tokenString string, secret string) (*entity.PayloadParams, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("token is not valid")
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		phone, ok := claims["phone"].(string)
		if !ok {
			return nil, errInvalidToken
		}
		refresh, ok := claims["refresh"].(bool)
		if !ok {
			return nil, errInvalidToken
		}
		return &entity.PayloadParams{
			Phone:   phone,
			Refresh: refresh,
		}, nil
	}
	return nil, errors.New("unable to parse claims")
}
