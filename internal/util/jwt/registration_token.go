package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	phoneKey             = "phone"
	code                 = "point_code"
	registrationTokenTTL = 1 * time.Hour
)

func GenerateRegistrationToken(phone, pointCode string, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		phoneKey: phone,
		code:     pointCode,
		"exp":    time.Now().Add(registrationTokenTTL).Unix(),
		"iat":    time.Now().Unix(),
	})
	return token.SignedString([]byte(secret))
}
