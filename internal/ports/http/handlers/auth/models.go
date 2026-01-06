package auth

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/command"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
)

//go:generate easyjson -all models.go

type sendCodeRequest struct {
	Phone string `json:"phone" validate:"required,min=10,max=15"`
}

func (s *sendCodeRequest) Validate() error {
	if err := request.GetValidator().Struct(s); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

type registerStaffRequest struct {
	Phone     string `json:"phone" validate:"required,min=10,max=15"`
	Role      string `json:"role" validate:"required,oneof=manager staff"`
	PointCode string `json:"point_code" validate:"required"`
}

func (r *registerStaffRequest) Validate() error {
	if err := request.GetValidator().Struct(r); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

func (r registerStaffRequest) toDomain() *command.RegisterStaffCommand {
	role, _ := enum.RoleString(r.Role)
	return &command.RegisterStaffCommand{
		Phone:     r.Phone,
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
	Name     string `json:"name" validate:"required,min=2"`
	Surname  string `json:"surname" validate:"required,min=2"`
	Password string `json:"password" validate:"required"`
}

func (r *registerConfirmRequest) Validate() error {
	if err := request.GetValidator().Struct(r); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

func (r *registerConfirmRequest) toDomain() *command.RegisterStaffConfirmCommand {
	return &command.RegisterStaffConfirmCommand{
		Token:    r.Token,
		Name:     r.Name,
		Surname:  r.Surname,
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

func (r *registerManagementRequest) toDomain() *command.RegisterManagementCommand {
	role, _ := enum.RoleString(r.Role)
	return &command.RegisterManagementCommand{
		Name:     r.Name,
		Surname:  r.Surname,
		Phone:    r.Phone,
		Role:     role,
		Password: r.Password,
		NetworkRegisterInfo: &command.NetworkRegisterInfo{
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
