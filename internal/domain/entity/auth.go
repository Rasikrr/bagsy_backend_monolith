package entity

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/golang-jwt/jwt"
)

type Auth struct {
	AccessToken  string
	RefreshToken string
}

type PayloadParams struct {
	jwt.StandardClaims

	Phone       string `json:"phone"`
	Role        string `json:"role"`
	Active      bool   `json:"active"`
	Refresh     bool   `json:"refresh"`
	PointCode   string `json:"point_code"`
	NetworkCode string `json:"network_code"`
}

func (p *PayloadParams) IsRefresh() bool {
	return p.Refresh
}

func (p *PayloadParams) GetRole() (enum.Role, error) {
	return enum.RoleString(p.Role)
}

func (p *PayloadParams) GetPhone() string {
	return p.Phone
}

func (p *PayloadParams) GetActive() bool {
	return p.Active
}

func (p *PayloadParams) ToSession() (*Session, error) {
	ses := NewSession().
		SetPhone(p.Phone).
		SetActive(p.Active).
		SetPointCode(p.PointCode)
	role, err := p.GetRole()
	if err != nil {
		return nil, err
	}
	ses.SetRole(role)
	return ses, nil
}
