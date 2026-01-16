package tokens

import authS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"

type inviteTokenInfoDTO struct {
	Phone       string `json:"phone"`
	PointCode   string `json:"point_code"`
	NetworkCode string `json:"network_code"`
	Purpose     string `json:"purpose"`
}

func (i *inviteTokenInfoDTO) toDomain() *authS.InviteTokenInfo {
	purpose, _ := authS.TokenPurposeString(i.Purpose)
	return &authS.InviteTokenInfo{
		Phone:       i.Phone,
		PointCode:   i.PointCode,
		NetworkCode: i.NetworkCode,
		Purpose:     purpose,
	}
}

func toDTO(info *authS.InviteTokenInfo) *inviteTokenInfoDTO {
	return &inviteTokenInfoDTO{
		Phone:       info.Phone,
		PointCode:   info.PointCode,
		NetworkCode: info.NetworkCode,
		Purpose:     info.Purpose.String(),
	}
}
