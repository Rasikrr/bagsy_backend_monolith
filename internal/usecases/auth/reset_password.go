package auth

import (
	"context"
	"fmt"
	"time"

	authDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/hasher"
	"github.com/cockroachdb/errors"
)

type resetTokenStore interface {
	Save(ctx context.Context, token string, phone shared.Phone, ttl time.Duration) error
	Get(ctx context.Context, token string) (shared.Phone, error)
	Delete(ctx context.Context, token string) error
}

type linkSender interface {
	SendPasswordResetLink(ctx context.Context, phone shared.Phone, link string) error
}

type ResetPasswordUseCase struct {
	employeeRepo   employeeRepository
	resetTokenRepo resetTokenStore
	tokenService   tokenService
	linkSender     linkSender
	resetTTL       time.Duration
	frontendURL    string
}

func NewResetPasswordUseCase(
	employeeRepo employeeRepository,
	resetTokenRepo resetTokenStore,
	tokenService tokenService,
	linkSender linkSender,
	resetTTL time.Duration,
	frontendURL string,
) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{
		employeeRepo:   employeeRepo,
		resetTokenRepo: resetTokenRepo,
		tokenService:   tokenService,
		linkSender:     linkSender,
		resetTTL:       resetTTL,
		frontendURL:    frontendURL,
	}
}

func (u *ResetPasswordUseCase) RequestReset(ctx context.Context, req RequestResetInput) error {
	phone, err := shared.NewPhone(req.Phone)
	if err != nil {
		return err
	}

	employee, err := u.employeeRepo.GetByPhone(ctx, phone)
	if err != nil {
		return errors.Wrap(err, "get employee by phone")
	}

	if !employee.IsActive() {
		return authDomain.ErrEmployeeInactive
	}

	resetToken, err := authDomain.NewResetToken(phone, u.resetTTL)
	if err != nil {
		return errors.Wrap(err, "generate reset token")
	}

	if err := u.resetTokenRepo.Save(ctx, resetToken.Token, phone, u.resetTTL); err != nil {
		return errors.Wrap(err, "save reset token")
	}

	link := fmt.Sprintf("%s/%s", u.frontendURL, resetToken.Token)

	if err := u.linkSender.SendPasswordResetLink(ctx, phone, link); err != nil {
		return errors.Wrap(err, "send password reset link")
	}

	return nil
}

func (u *ResetPasswordUseCase) ConfirmReset(ctx context.Context, req ConfirmResetInput) error {
	phone, err := u.resetTokenRepo.Get(ctx, req.Token)
	if err != nil {
		return errors.Wrap(err, "get reset token")
	}

	employee, err := u.employeeRepo.GetByPhone(ctx, phone)
	if err != nil {
		return errors.Wrap(err, "get employee for password reset")
	}

	passwordHash, err := hasher.Password(req.NewPassword)
	if err != nil {
		return errors.Wrap(err, "hash new password")
	}

	if err := employee.ChangePassword(passwordHash); err != nil {
		return err
	}

	if err := u.employeeRepo.Save(ctx, employee); err != nil {
		return errors.Wrap(err, "save employee after password change")
	}

	if err := u.tokenService.DeleteAllRefreshTokens(ctx, employee.ID); err != nil {
		return errors.Wrap(err, "invalidate all sessions")
	}

	_ = u.resetTokenRepo.Delete(ctx, req.Token)

	return nil
}
