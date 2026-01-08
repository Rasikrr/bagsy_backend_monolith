package jwt

import (
	"github.com/golang-jwt/jwt"
)

// claims для access токена с полным набором данных
type accessClaims struct {
	jwt.StandardClaims

	Phone       string `json:"phone"`
	Role        string `json:"role,omitempty"`
	PointCode   string `json:"point_code,omitempty"`
	NetworkCode string `json:"network_code,omitempty"`
}
