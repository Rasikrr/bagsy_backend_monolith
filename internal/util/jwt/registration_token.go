package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	registrationTokenTTL = 24 * time.Hour
)

func GenerateRegistrationToken(phone, pointCode, networkCode string, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		phoneKey:   phone,
		pointKey:   pointCode,
		networkKey: networkCode,
		"exp":      time.Now().Add(registrationTokenTTL).Unix(),
		"iat":      time.Now().Unix(),
	})
	return token.SignedString([]byte(secret))
}

// nolint: govet
func ValidateRegistrationToken(tokenString, secret string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", errUnexpectedSigning, token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if _, ok := claims[phoneKey]; !ok {
			return errInvalidToken
		}
		if _, ok := claims[pointKey]; !ok {
			return errInvalidToken
		}
		if _, ok := claims[networkKey]; !ok {
			return errInvalidToken
		}
		return nil
	}
	return errInvalidToken
}
