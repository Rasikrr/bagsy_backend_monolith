package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/domain/entity"
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

func GenerateRefreshToken(params *entity.PayloadParams) (string, error) {
	claims := jwt.MapClaims{
		"phone":   params.Phone,
		"role":    params.Role,
		"refresh": params.Refresh,
		"exp":     time.Now().Add(refreshTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
