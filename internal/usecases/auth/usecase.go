package auth

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/hasher"
)

type UseCase struct {
}

func NewUseCase() *UseCase {
	return &UseCase{}
}

func (u *UseCase) LoginEmployee(ctx context.Context, phone, password string) (*auth.Token, error) {
	phoneVO, err := shared.NewPhone(phone)
	if err != nil {
		return nil, err
	}
	hasher.CheckPassword()

}
