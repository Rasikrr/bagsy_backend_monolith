package jwt

import (
	"github.com/golang-jwt/jwt"
)

type accessClaims struct {
	jwt.StandardClaims

	Phone  string `json:"phone"`
	UserID string `json:"user_id"`
}
