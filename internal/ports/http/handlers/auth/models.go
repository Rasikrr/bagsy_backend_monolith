package auth

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	authS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"
)

//go:generate easyjson -all models.go

type changePasswordRequest struct {
	Phone string `json:"phone" validate:"required,min=10,max=15"`
}

func (c *changePasswordRequest) Validate() error {
	if err := request.GetValidator().Struct(c); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

type registerStaffRequest struct {
	Phone     string `json:"phone" validate:"required,min=10,max=15"`
	Name      string `json:"name" validate:"required"`
	Surname   string `json:"surname" validate:"required"`
	Role      string `json:"role" validate:"required,oneof=manager staff"`
	PointCode string `json:"point_code" validate:"required"`
}

func (r *registerStaffRequest) Validate() error {
	if err := request.GetValidator().Struct(r); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

func (r registerStaffRequest) toDomain() *auth.RegisterStaffCommand {
	role, _ := user.RoleString(r.Role)
	return &auth.RegisterStaffCommand{
		Phone:     r.Phone,
		Name:      r.Name,
		Surname:   r.Surname,
		Role:      role,
		PointCode: r.PointCode,
	}
}

type loginRequest struct {
	Phone    string `json:"phone"    validate:"required,min=10,max=15"`
	Password string `json:"password" validate:"required"`
}

func (l *loginRequest) Validate() error {
	if err := request.GetValidator().Struct(l); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type refreshTokensRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r *refreshTokensRequest) Validate() error {
	if err := request.GetValidator().Struct(r); err != nil {
		return request.HandleValidationError(err)
	}

	if r.RefreshToken == "" {
		return domainErr.NewInvalidInputError("refresh_token is required", nil)
	}

	return nil
}

type refreshTokensResponse loginResponse

type registerConfirmRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type passwordChangeConfirmRequest registerConfirmRequest

func (p *passwordChangeConfirmRequest) Validate() error {
	if err := request.GetValidator().Struct(p); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

func (p *passwordChangeConfirmRequest) toDomain() *auth.ChangePasswordConfirmCommand {
	return &auth.ChangePasswordConfirmCommand{
		Token:    p.Token,
		Password: p.Password,
	}
}

func (r *registerConfirmRequest) Validate() error {
	if err := request.GetValidator().Struct(r); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

func (r *registerConfirmRequest) toDomain() *auth.RegisterStaffConfirmCommand {
	return &auth.RegisterStaffConfirmCommand{
		Token:    r.Token,
		Password: r.Password,
	}
}

type registerConfirmResponse loginResponse

type networkInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type registerManagementRequest struct {
	Name        string      `json:"name" validate:"required,min=2"`
	Surname     string      `json:"surname" validate:"required,min=2"`
	Phone       string      `json:"phone" validate:"required"`
	Password    string      `json:"password" validate:"required"`
	Role        string      `json:"role" validate:"required,oneof=net_manager self_owner"`
	NetworkInfo networkInfo `json:"network_info" validate:"required"`
}

func (r *registerManagementRequest) Validate() error {
	if err := request.GetValidator().Struct(r); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

func (r *registerManagementRequest) toDomain() *auth.RegisterManagementCommand {
	role, _ := user.RoleString(r.Role)
	return &auth.RegisterManagementCommand{
		Name:     r.Name,
		Surname:  r.Surname,
		Phone:    r.Phone,
		Role:     role,
		Password: r.Password,
		NetworkRegisterInfo: &auth.RegisterNetworkInfo{
			Name:        r.NetworkInfo.Name,
			Description: r.NetworkInfo.Description,
		},
	}
}

type registerManagementConfirmRequest struct {
	Phone string `json:"phone" validate:"required"`
	Code  string `json:"code" validate:"required"`
}

func (r *registerManagementConfirmRequest) Validate() error {
	if err := request.GetValidator().Struct(r); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

type resendRegisterManagementRequest struct {
	Phone string `json:"phone" validate:"required"`
}

func (r *resendRegisterManagementRequest) Validate() error {
	if err := request.GetValidator().Struct(r); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

type registerStaffResendRequest = changePasswordRequest

type verifyAuthTokenResponse struct {
	Phone       string `json:"phone"`
	NetworkCode string `json:"network_code,omitempty"`
	PointCode   string `json:"point_code,omitempty"`
	Purpose     string `json:"purpose"`
}

func toVerifyAuthTokenResponse(dto *authS.InviteTokenInfo) *verifyAuthTokenResponse {
	return &verifyAuthTokenResponse{
		Phone:       dto.Phone,
		NetworkCode: dto.NetworkCode,
		PointCode:   dto.PointCode,
		Purpose:     dto.Purpose.String(),
	}
}
