package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	phoneKey             = "phone"
	registrationTokenTTL = 1 * time.Hour
)

func GenerateRegistrationToken(phone string, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		phoneKey: phone,
		"exp":    time.Now().Add(registrationTokenTTL).Unix(),
		"iat":    time.Now().Unix(),
	})
	return token.SignedString([]byte(secret))
}
