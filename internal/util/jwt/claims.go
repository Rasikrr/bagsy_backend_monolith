package jwt

import (
	"errors"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/deref"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims

	Phone       string `json:"phone"`
	Role        string `json:"role"`
	Refresh     bool   `json:"refresh"`
	PointCode   string `json:"point_code"`
	NetworkCode string `json:"network_code"`
}

func NewClaims(user *entity.User, ttl time.Duration, isRefresh bool) Claims {
	c := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "system",
		},
		Phone:   user.Phone,
		Refresh: isRefresh,
	}
	if !isRefresh {
		c.Role = user.Role.String()
		c.PointCode = deref.String(user.PointCode)
		c.NetworkCode = deref.String(user.NetworkCode)
	}
	return c
}

// ToSession конвертирует Claims в доменную сущность Session
func (c *Claims) ToSession() (*entity.Session, error) {
	if c.Refresh {
		return nil, errors.New("refresh token cannot be converted to session")
	}

	roleEnum, err := enum.RoleString(c.Role)
	if err != nil {
		return nil, err
	}

	session := entity.NewSession().
		SetPhone(c.Phone).
		SetRole(roleEnum)

	if c.PointCode != "" {
		session.SetPointCode(c.PointCode)
	}

	if c.NetworkCode != "" {
		session.SetNetworkCode(c.NetworkCode)
	}

	return session, nil
}
