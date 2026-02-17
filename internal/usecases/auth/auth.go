package auth

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/hasher"
	"github.com/google/uuid"
)

type employeeGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*identity.Employee, error)
}

type UseCase struct {
	employeeRepository employeeRepository
	employeeGetter     employeeGetter
	tokenService       tokenService
}

func NewUseCase(
	employeeRepository employeeRepository,
	employeeGetter employeeGetter,
	tokenService tokenService,
) *UseCase {
	return &UseCase{
		employeeRepository: employeeRepository,
		employeeGetter:     employeeGetter,
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
	if !hasher.CheckPassword(employee.PasswordHash, password) {
		return nil, identity.ErrEmployeeNotFound
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

func (u *UseCase) RefreshTokens(ctx context.Context, refreshToken string) (*TokensOutput, error) {
	// 1. Validate + delete old refresh token, get userID.
	userID, err := u.tokenService.RefreshTokens(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("refresh tokens: %w", err)
	}

	// 2. Load employee to get phone for new access token claims.
	employee, err := u.employeeGetter.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get employee for refresh: %w", err)
	}

	// 3. Generate new token pair.
	access, refresh, err := u.tokenService.GenerateTokens(ctx, employee.ID, employee.Phone)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return &TokensOutput{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (u *UseCase) Logout(ctx context.Context, refreshToken string) error {
	return u.tokenService.DeleteRefreshToken(ctx, refreshToken)
}
