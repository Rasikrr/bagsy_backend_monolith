package employee

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
	"github.com/google/uuid"
)

type employeeRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*identity.Employee, error)
}

type mediaRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*media.Asset, error)
}

type storageClient interface {
	GeneratePresignedDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (string, error)
}

type UseCase struct {
	employeeRepo    employeeRepository
	mediaRepo       mediaRepository
	storage         storageClient
	avatarURLExpiry time.Duration
}

func NewUseCase(
	employeeRepo employeeRepository,
	mediaRepo mediaRepository,
	storage storageClient,
	avatarURLExpiry time.Duration,
) *UseCase {
	return &UseCase{
		employeeRepo:    employeeRepo,
		mediaRepo:       mediaRepo,
		storage:         storage,
		avatarURLExpiry: avatarURLExpiry,
	}
}
