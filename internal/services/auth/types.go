package auth

import "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"

type ManagementRegistrationState struct {
	Command  *auth.RegisterManagementCommand
	AuthCode string
	Attempts int
}

func newManagementRegistrationState(cmd *auth.RegisterManagementCommand, authCode string) *ManagementRegistrationState {
	return &ManagementRegistrationState{
		Command:  cmd,
		AuthCode: authCode,
	}
}

type InviteTokenInfo struct {
	Phone       string
	PointCode   string
	NetworkCode string
	Purpose     TokenPurpose
}

type AccessTokenPayload struct {
	Phone       string
	Role        string
	PointCode   string
	NetworkCode string
}

func newAccessTokenPayload(phone, role, pointCode, networkCode string) *AccessTokenPayload {
	return &AccessTokenPayload{
		Phone:       phone,
		Role:        role,
		PointCode:   pointCode,
		NetworkCode: networkCode,
	}
}
