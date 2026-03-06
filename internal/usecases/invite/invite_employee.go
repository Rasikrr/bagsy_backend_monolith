package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	authDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/hasher"
	"github.com/Rasikrr/core/log"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

const (
	inviteCooldown = 60 * time.Second
)

type employeeRepository interface {
	CountByOrganization(ctx context.Context, orgID uuid.UUID) (int, error)
	GetByPhone(ctx context.Context, phone shared.Phone) (*identity.Employee, error)
	ExistsByPhone(ctx context.Context, phone shared.Phone) (bool, error)
	Save(ctx context.Context, emp *identity.Employee) error
}

type workHistoryRepository interface {
	Save(ctx context.Context, wh *identity.WorkHistory) error
}

type actionTokenStore interface {
	Save(ctx context.Context, token *authDomain.ActionToken) error
	Get(ctx context.Context, token string) (*authDomain.ActionToken, error)
	Delete(ctx context.Context, token string) error
}

type pendingInviteStore interface {
	Save(ctx context.Context, inv *PendingInvite) error
	Get(ctx context.Context, phone shared.Phone) (*PendingInvite, error)
	Delete(ctx context.Context, phone shared.Phone) error
}

type inviteLinkSender interface {
	SendInviteLink(ctx context.Context, phone shared.Phone, link string) error
}

type tokenService interface {
	GenerateTokens(ctx context.Context, userID uuid.UUID, phone shared.Phone) (access, refresh string, err error)
}

type invitePolicy interface {
	CanInviteEmployee(orgCtx *access.OrgContext, targetRole identity.Role, currentCount int) error
}

type txManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type UseCase struct {
	employeeRepo    employeeRepository
	workHistoryRepo workHistoryRepository
	actionTokenRepo actionTokenStore
	pendingInvRepo  pendingInviteStore
	tokenService    tokenService
	linkSender      inviteLinkSender
	policy          invitePolicy
	txManager       txManager
	inviteTTL       time.Duration
	frontendURL     string
}

func NewUseCase(
	employeeRepo employeeRepository,
	workHistoryRepo workHistoryRepository,
	actionTokenRepo actionTokenStore,
	pendingInvRepo pendingInviteStore,
	tokenService tokenService,
	linkSender inviteLinkSender,
	policy invitePolicy,
	txManager txManager,
	inviteTTL time.Duration,
	frontendURL string,
) *UseCase {
	return &UseCase{
		employeeRepo:    employeeRepo,
		workHistoryRepo: workHistoryRepo,
		actionTokenRepo: actionTokenRepo,
		pendingInvRepo:  pendingInvRepo,
		tokenService:    tokenService,
		linkSender:      linkSender,
		policy:          policy,
		txManager:       txManager,
		inviteTTL:       inviteTTL,
		frontendURL:     frontendURL,
	}
}

func (u *UseCase) SendInvite(ctx context.Context, orgCtx *access.OrgContext, input SendInviteInput) (*SendInviteOutput, error) {
	phone, err := shared.NewPhone(input.Phone)
	if err != nil {
		return nil, err
	}

	role, err := identity.ParseRole(input.Role)
	if err != nil {
		return nil, err
	}

	employeesCount, err := u.employeeRepo.CountByOrganization(ctx, orgCtx.Organization.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to count employees by organization: %w", err)
	}

	if err = u.policy.CanInviteEmployee(orgCtx, role, employeesCount); err != nil {
		return nil, err
	}

	exists, err := u.employeeRepo.ExistsByPhone(ctx, phone)
	if err != nil {
		return nil, errors.Wrap(err, "check phone uniqueness")
	}
	if exists {
		return nil, authDomain.ErrPhoneAlreadyExists
	}

	if existing, _ := u.pendingInvRepo.Get(ctx, phone); existing != nil {
		if time.Since(existing.LastSentAt) < inviteCooldown {
			return nil, authDomain.ErrInviteAlreadyExists
		}
	}

	locationID := &orgCtx.Employee.LocationID
	inviteToken, err := authDomain.NewStaffInviteToken(
		phone,
		locationID,
		orgCtx.Organization.ID,
		u.inviteTTL,
	)
	if err != nil {
		return nil, errors.Wrap(err, "generate invite token")
	}

	permissions := identity.DefaultPermissionsForRole(role)
	now := time.Now()

	pending := &PendingInvite{
		Phone:          phone,
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		OrganizationID: orgCtx.Organization.ID,
		LocationID:     input.LocationID,
		Role:           role,
		Permissions:    permissions,
		InvitedBy:      orgCtx.Employee.ID,
		LastSentAt:     now,
		ExpiresAt:      now.Add(u.inviteTTL),
	}

	if err = u.pendingInvRepo.Save(ctx, pending); err != nil {
		return nil, errors.Wrap(err, "save pending invite")
	}

	if err = u.actionTokenRepo.Save(ctx, inviteToken); err != nil {
		return nil, errors.Wrap(err, "save invite token")
	}

	link := fmt.Sprintf("%s/%s", u.frontendURL, inviteToken.Token)
	if err = u.linkSender.SendInviteLink(ctx, phone, link); err != nil {
		return nil, errors.Wrap(err, "send invite link")
	}
	log.Infof(ctx, "invite link sent to %s", link)

	return &SendInviteOutput{
		Phone:     phone.String(),
		ExpiresIn: int(u.inviteTTL.Seconds()),
	}, nil
}

func (u *UseCase) ConfirmInvite(ctx context.Context, input ConfirmInviteInput) (*TokensOutput, error) {
	actionToken, err := u.actionTokenRepo.Get(ctx, input.Token)
	if err != nil {
		return nil, errors.Wrap(err, "get invite token")
	}

	phone := actionToken.Phone

	pending, err := u.pendingInvRepo.Get(ctx, phone)
	if err != nil {
		return nil, errors.Wrap(err, "get pending invite")
	}
	if pending == nil {
		return nil, authDomain.ErrInviteTokenExpired
	}

	passwordHash, err := hasher.Password(input.Password)
	if err != nil {
		return nil, errors.Wrap(err, "hash password")
	}

	var employeeID uuid.UUID

	err = u.txManager.Do(ctx, func(txCtx context.Context) error {
		exists, e := u.employeeRepo.ExistsByPhone(txCtx, phone)
		if e != nil {
			return errors.Wrap(e, "check phone uniqueness")
		}
		if exists {
			return authDomain.ErrPhoneAlreadyExists
		}

		emp, e := identity.NewEmployee(identity.CreateEmployeeParams{
			Phone:          phone,
			FirstName:      pending.FirstName,
			LastName:       pending.LastName,
			OrganizationID: pending.OrganizationID,
			LocationID:     pending.LocationID,
			Role:           pending.Role,
			Permissions:    pending.Permissions,
		})
		if e != nil {
			return errors.Wrap(e, "create employee")
		}
		emp.SetPassword(passwordHash)

		if e = u.employeeRepo.Save(txCtx, emp); e != nil {
			return errors.Wrap(e, "save employee")
		}
		employeeID = emp.ID

		comment := "Приглашён в организацию"
		wh := identity.NewWorkHistory(
			employeeID,
			pending.OrganizationID,
			pending.LocationID,
			pending.Role,
			identity.ChangeTypeHired,
			&comment,
		)
		if e = u.workHistoryRepo.Save(txCtx, wh); e != nil {
			return errors.Wrap(e, "save work history")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	_ = u.pendingInvRepo.Delete(ctx, phone)
	_ = u.actionTokenRepo.Delete(ctx, input.Token)

	access, refresh, err := u.tokenService.GenerateTokens(ctx, employeeID, phone)
	if err != nil {
		return nil, errors.Wrap(err, "generate tokens")
	}

	return &TokensOutput{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (u *UseCase) ResendInvite(ctx context.Context, orgCtx *access.OrgContext, input ResendInviteInput) (*ResendInviteOutput, error) {
	phone, err := shared.NewPhone(input.Phone)
	if err != nil {
		return nil, err
	}

	pending, err := u.pendingInvRepo.Get(ctx, phone)
	if err != nil {
		return nil, errors.Wrap(err, "get pending invite")
	}
	if pending == nil {
		return nil, authDomain.ErrInviteTokenNotFound
	}

	if pending.OrganizationID != orgCtx.Organization.ID {
		return nil, identity.ErrPermissionDenied
	}

	if time.Since(pending.LastSentAt) < inviteCooldown {
		return nil, authDomain.ErrInviteAlreadyExists
	}

	inviteToken, err := authDomain.NewStaffInviteToken(
		phone,
		pending.LocationID,
		pending.OrganizationID,
		u.inviteTTL,
	)
	if err != nil {
		return nil, errors.Wrap(err, "generate invite token")
	}

	now := time.Now()
	pending.LastSentAt = now
	pending.ExpiresAt = now.Add(u.inviteTTL)

	if err = u.pendingInvRepo.Save(ctx, pending); err != nil {
		return nil, errors.Wrap(err, "save pending invite")
	}

	if err = u.actionTokenRepo.Save(ctx, inviteToken); err != nil {
		return nil, errors.Wrap(err, "save invite token")
	}

	link := fmt.Sprintf("%s/%s", u.frontendURL, inviteToken.Token)
	if err = u.linkSender.SendInviteLink(ctx, phone, link); err != nil {
		return nil, errors.Wrap(err, "send invite link")
	}

	return &ResendInviteOutput{
		Phone:      phone.String(),
		ExpiresIn:  int(u.inviteTTL.Seconds()),
		RetryAfter: int(inviteCooldown.Seconds()),
	}, nil
}
