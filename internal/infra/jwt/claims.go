package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type accessClaims struct {
	jwt.RegisteredClaims

	Phone  string `json:"phone"`
	UserID string `json:"user_id"`
}
