package tokens

import authS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"

type inviteTokenInfoDTO struct {
	Phone       string `json:"phone"`
	PointCode   string `json:"point_code"`
	NetworkCode string `json:"network_code"`
}

func (i *inviteTokenInfoDTO) toDomain() *authS.InviteTokenInfo {
	return &authS.InviteTokenInfo{
		Phone:       i.Phone,
		PointCode:   i.PointCode,
		NetworkCode: i.NetworkCode,
	}
}

func toDTO(info *authS.InviteTokenInfo) *inviteTokenInfoDTO {
	return &inviteTokenInfoDTO{
		Phone:       info.Phone,
		PointCode:   info.PointCode,
		NetworkCode: info.NetworkCode,
	}
}
