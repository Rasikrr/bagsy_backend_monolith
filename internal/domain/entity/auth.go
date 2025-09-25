package entity

import "github.com/Rasikrr/bugsy_backend_monolith/internal/domain/enum"

type Auth struct {
	AccessToken  string
	RefreshToken string
}

type PayloadParams struct {
	Phone     string
	Role      string
	Active    bool
	Refresh   bool
	PointCode string
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
