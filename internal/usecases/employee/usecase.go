package employee

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type employeeRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*identity.Employee, error)
	GetByFilter(ctx context.Context, filter *identity.EmployeeFilter) (*identity.EmployeePage, error)
}

type mediaRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*media.Asset, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*media.Asset, error)
}

type storageClient interface {
	GeneratePresignedDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (string, error)
}

type policy interface {
	CanListEmployees(orgCtx *access.OrgContext, filter *identity.EmployeeFilter) error
}

type UseCase struct {
	employeeRepo    employeeRepository
	mediaRepo       mediaRepository
	storage         storageClient
	policy          policy
	avatarURLExpiry time.Duration
}

func NewUseCase(
	employeeRepo employeeRepository,
	mediaRepo mediaRepository,
	storage storageClient,
	policy policy,
	avatarURLExpiry time.Duration,
) *UseCase {
	return &UseCase{
		employeeRepo:    employeeRepo,
		mediaRepo:       mediaRepo,
		storage:         storage,
		policy:          policy,
		avatarURLExpiry: avatarURLExpiry,
	}
}

// nolint: nilnil
func (u *UseCase) resolveAvatarURL(ctx context.Context, avatarID *uuid.UUID) (*string, error) {
	if avatarID == nil {
		return nil, nil
	}

	asset, err := u.mediaRepo.GetByID(ctx, *avatarID)
	if err != nil {
		return nil, errors.Wrap(err, "get avatar asset")
	}

	if !asset.IsReady() {
		return nil, nil
	}

	url, err := u.storage.GeneratePresignedDownloadURL(ctx, asset.ObjectKey, u.avatarURLExpiry)
	if err != nil {
		return nil, errors.Wrap(err, "generate avatar url")
	}

	return &url, nil
}

// resolveAvatarURLsBatch загружает все ассеты одним запросом и генерирует presigned URLs.
// Возвращает map[avatarID] → presigned URL.
func (u *UseCase) resolveAvatarURLsBatch(ctx context.Context, avatarIDs []uuid.UUID) (map[uuid.UUID]string, error) {
	if len(avatarIDs) == 0 {
		return map[uuid.UUID]string{}, nil
	}

	assets, err := u.mediaRepo.GetByIDs(ctx, avatarIDs)
	if err != nil {
		return nil, errors.Wrap(err, "get avatar assets")
	}

	result := make(map[uuid.UUID]string, len(assets))
	for _, asset := range assets {
		if !asset.IsReady() {
			continue
		}
		var url string
		url, err = u.storage.GeneratePresignedDownloadURL(ctx, asset.ObjectKey, u.avatarURLExpiry)
		if err != nil {
			return nil, errors.Wrap(err, "generate avatar url")
		}
		result[asset.ID] = url
	}

	return result, nil
}
