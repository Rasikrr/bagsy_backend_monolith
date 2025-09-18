package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	phoneKey             = "phone"
	registrationTokenTTL = 1 * time.Hour
)

func GenerateRegistrationToken(phone string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		phoneKey: phone,
		"exp":    time.Now().Add(registrationTokenTTL).Unix(),
		"iat":    time.Now().Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
