package auth

import (
	"context"
	"time"

	authDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/organization"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/hasher"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

const (
	otpTTL         = 5 * time.Minute
	otpCooldown    = 60 * time.Second
	otpMaxAttempts = 3
)

type employeeRepository interface {
	GetByPhone(ctx context.Context, phone shared.Phone) (*identity.Employee, error)
	ExistsByPhone(ctx context.Context, phone shared.Phone) (bool, error)
	Save(ctx context.Context, emp *identity.Employee) error
}

type organizationRepository interface {
	Save(ctx context.Context, org *organization.Organization) error
}

type planRepository interface {
	FindActiveByCode(ctx context.Context, code billing.PlanCode) (*billing.Plan, error)
}

type subscriptionRepository interface {
	Save(ctx context.Context, sub *billing.Subscription) error
}

type workHistoryRepository interface {
	Save(ctx context.Context, wh *identity.WorkHistory) error
}

type pendingRegistrationStore interface {
	Save(ctx context.Context, reg *PendingRegistration) error
	Get(ctx context.Context, phone shared.Phone) (*PendingRegistration, error)
	Delete(ctx context.Context, phone shared.Phone) error
}

type otpSender interface {
	SendOTP(ctx context.Context, phone shared.Phone, code string) error
}

type tokenService interface {
	GenerateTokens(ctx context.Context, userID uuid.UUID, phone shared.Phone) (access, refresh string, err error)
	RefreshTokens(ctx context.Context, oldRefreshToken string) (userID uuid.UUID, err error)
	DeleteRefreshToken(ctx context.Context, refresh string) error
	DeleteAllRefreshTokens(ctx context.Context, userID uuid.UUID) error
}

type txManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type RegisterOwnerUseCase struct {
	employeesRepo        employeeRepository
	plansRepo            planRepository
	pendingRequestsStore pendingRegistrationStore
	organizationRepo     organizationRepository
	subscriptionsRepo    subscriptionRepository
	workHistoryRepo      workHistoryRepository
	tokensService        tokenService
	otpSender            otpSender
	txManager            txManager
}

func NewRegisterOwnerUseCase(
	employeesRepo employeeRepository,
	plansRepo planRepository,
	organizationRepo organizationRepository,
	subscriptions subscriptionRepository,
	workHistory workHistoryRepository,
	tokenService tokenService,
	pendingRequestsStore pendingRegistrationStore,
	txManager txManager,
	otpSender otpSender,
) *RegisterOwnerUseCase {
	return &RegisterOwnerUseCase{
		employeesRepo:        employeesRepo,
		plansRepo:            plansRepo,
		organizationRepo:     organizationRepo,
		subscriptionsRepo:    subscriptions,
		pendingRequestsStore: pendingRequestsStore,
		workHistoryRepo:      workHistory,
		otpSender:            otpSender,
		tokensService:        tokenService,
		txManager:            txManager,
	}
}

func (u *RegisterOwnerUseCase) Register(ctx context.Context, req RegisterInput) (*RegisterOutput, error) {
	phone, err := shared.NewPhone(req.Phone)
	if err != nil {
		return nil, err
	}

	exists, err := u.employeesRepo.ExistsByPhone(ctx, phone)
	if err != nil {
		return nil, errors.Wrap(err, "check phone uniqueness")
	}
	if exists {
		return nil, authDomain.ErrPhoneAlreadyExists
	}

	planCode, err := billing.ParsePlanCode(req.PlanCode)
	if err != nil {
		return nil, err
	}

	if _, err := u.plansRepo.FindActiveByCode(ctx, planCode); err != nil {
		return nil, err
	}

	// If a pending registration already exists for this phone,
	// treat it as a re-submit: overwrite data, send new code
	// (only if cooldown has passed).
	if existing, _ := u.pendingRequestsStore.Get(ctx, phone); existing != nil {
		if time.Since(existing.LastSentAt) < otpCooldown {
			return nil, authDomain.ErrOTPAlreadySent
		}
	}

	passwordHash, err := hasher.Password(req.Password)
	if err != nil {
		return nil, errors.Wrap(err, "hash password")
	}

	otp, err := authDomain.NewOTPCode(phone, otpTTL)
	if err != nil {
		return nil, errors.Wrap(err, "generate OTP")
	}

	now := time.Now()

	pending := &PendingRegistration{
		Phone:        phone,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PasswordHash: passwordHash,
		PlanCode:     planCode,
		OTPCode:      otp.Code,
		Attempts:     0,
		MaxAttempts:  otpMaxAttempts,
		LastSentAt:   now,
		ExpiresAt:    now.Add(otpTTL),
	}

	if err := u.pendingRequestsStore.Save(ctx, pending); err != nil {
		return nil, errors.Wrap(err, "save pending registration")
	}

	if err := u.otpSender.SendOTP(ctx, phone, otp.Code); err != nil {
		return nil, errors.Wrap(err, "send OTP")
	}

	return &RegisterOutput{
		Phone:      phone.String(),
		ExpiresIn:  int(otpTTL.Seconds()),
		RetryAfter: int(otpCooldown.Seconds()),
	}, nil
}

func (u *RegisterOwnerUseCase) VerifyRegistration(ctx context.Context, req VerifyInput) (*TokensOutput, error) {
	phone, err := shared.NewPhone(req.Phone)
	if err != nil {
		return nil, err
	}

	reg, err := u.pendingRequestsStore.Get(ctx, phone)
	if err != nil {
		return nil, errors.Wrap(err, "get pending registration")
	}
	if reg == nil {
		return nil, authDomain.ErrRegistrationExpired
	}

	if reg.Attempts >= reg.MaxAttempts {
		return nil, authDomain.ErrTooManyAttempts
	}

	if reg.OTPCode != req.Code {
		reg.Attempts++
		_ = u.pendingRequestsStore.Save(ctx, reg)
		if reg.Attempts >= reg.MaxAttempts {
			return nil, authDomain.ErrTooManyAttempts
		}
		return nil, authDomain.ErrOTPInvalid
	}

	plan, err := u.plansRepo.FindActiveByCode(ctx, reg.PlanCode)
	if err != nil {
		return nil, errors.Wrap(err, "find plan")
	}

	var (
		employeeID uuid.UUID
		orgID      uuid.UUID
	)

	err = u.txManager.Do(ctx, func(txCtx context.Context) error {
		// Race condition protection: re-check phone uniqueness inside tx.
		exists, err := u.employeesRepo.ExistsByPhone(txCtx, phone)
		if err != nil {
			return errors.Wrap(err, "check phone uniqueness")
		}
		if exists {
			return authDomain.ErrPhoneAlreadyExists
		}

		// 1. Create stub organization.
		org, err := organization.NewStubOrganization()
		if err != nil {
			return errors.Wrap(err, "create organization")
		}
		if err := u.organizationRepo.Save(txCtx, org); err != nil {
			return errors.Wrap(err, "save organization")
		}
		orgID = org.ID

		employeePermissions := identity.NewPermissions(false, true)

		if plan.Code.IsSolo() {
			employeePermissions = identity.NewPermissions(true, true)
		}

		// 2. Create employee (owner).
		emp, err := identity.NewOwnerEmployee(identity.CreateOwnerParams{
			Phone:          phone,
			FirstName:      reg.FirstName,
			LastName:       reg.LastName,
			OrganizationID: orgID,
			Permissions:    employeePermissions,
		})
		if err != nil {
			return errors.Wrap(err, "create employee")
		}
		emp.SetPassword(reg.PasswordHash)

		if err := u.employeesRepo.Save(txCtx, emp); err != nil {
			return errors.Wrap(err, "save employee")
		}
		employeeID = emp.ID

		// 3. Create trial subscription.
		sub := billing.NewTrialSubscription(orgID, plan.ID, billing.DefaultTrialDays)
		if err := u.subscriptionsRepo.Save(txCtx, sub); err != nil {
			return errors.Wrap(err, "save subscription")
		}

		// 4. Create work history entry.
		comment := "Организация создана"
		wh := identity.NewWorkHistory(
			employeeID,
			orgID,
			nil,
			identity.RoleOwner,
			identity.ChangeTypeHired,
			&comment,
		)
		if err := u.workHistoryRepo.Save(txCtx, wh); err != nil {
			return errors.Wrap(err, "save work history")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Clean up pending registration after successful commit.
	_ = u.pendingRequestsStore.Delete(ctx, phone)

	// Generate JWT tokens.
	access, refresh, err := u.tokensService.GenerateTokens(ctx, employeeID, phone)
	if err != nil {
		return nil, err
	}

	return &TokensOutput{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (u *RegisterOwnerUseCase) Resend(ctx context.Context, req ResendInput) (*ResendOutput, error) {
	phone, err := shared.NewPhone(req.Phone)
	if err != nil {
		return nil, err
	}

	reg, err := u.pendingRequestsStore.Get(ctx, phone)
	if err != nil {
		return nil, errors.Wrap(err, "get pending registration")
	}
	if reg == nil {
		return nil, authDomain.ErrRegistrationExpired
	}

	if time.Since(reg.LastSentAt) < otpCooldown {
		return nil, authDomain.ErrOTPAlreadySent
	}

	otp, err := authDomain.NewOTPCode(phone, otpTTL)
	if err != nil {
		return nil, errors.Wrap(err, "generate OTP")
	}

	now := time.Now()
	reg.OTPCode = otp.Code
	reg.Attempts = 0
	reg.LastSentAt = now
	reg.ExpiresAt = now.Add(otpTTL)

	if err := u.pendingRequestsStore.Save(ctx, reg); err != nil {
		return nil, errors.Wrap(err, "save pending registration")
	}

	if err := u.otpSender.SendOTP(ctx, phone, otp.Code); err != nil {
		return nil, errors.Wrap(err, "send OTP")
	}

	return &ResendOutput{
		ExpiresIn:  int(otpTTL.Seconds()),
		RetryAfter: int(otpCooldown.Seconds()),
	}, nil
}
