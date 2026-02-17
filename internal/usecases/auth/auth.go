package auth

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/hasher"
)

type UseCase struct {
	employeeRepository employeeRepository
	tokenService       tokenService
}

func NewUseCase(employeeRepository employeeRepository, tokenService tokenService) *UseCase {
	return &UseCase{
		employeeRepository: employeeRepository,
		tokenService:       tokenService,
	}
}

func (u *UseCase) LoginEmployee(ctx context.Context, phone, password string) (*TokensOutput, error) {
	phoneVo, err := shared.NewPhone(phone)
	if err != nil {
		return nil, err
	}
	employee, err := u.employeeRepository.GetByPhone(ctx, phoneVo)
	if err != nil {
		return nil, err
	}
	ok := hasher.CheckPassword(employee.PasswordHash, password)
	if !ok {
		return nil, err
	}
	access, refresh, err := u.tokenService.GenerateTokens(ctx, employee.ID, employee.Phone)
	if err != nil {
		return nil, err
	}
	return &TokensOutput{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (u *UseCase) Logout(ctx context.Context, refreshToken string) error {
	return u.tokenService.DeleteRefreshToken(ctx, refreshToken)
}
