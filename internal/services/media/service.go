package media

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/dto"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/session"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/util/ptr"
	"github.com/Rasikrr/core/database"
	"github.com/google/uuid"
)

type storage interface {
	GetBucketName(ctx context.Context) string
	GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, expiresIn time.Duration) (string, error)
	GeneratePresignedDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (string, error)
}

type mediaRepository interface {
	Create(ctx context.Context, media *entity.Media) error
	GetByID(ctx context.Context, mediaID uuid.UUID) (*entity.Media, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status enum.MediaStatus) error
	SoftDeleteByID(ctx context.Context, id uuid.UUID) error
}

type userAvatarRepository interface {
	Get(ctx context.Context, phone string) (*entity.UserMedia, error)
	Set(ctx context.Context, userMedia *entity.UserMedia) error
	GetWithMedia(ctx context.Context, phone string) (*entity.Media, error)
	Remove(ctx context.Context, phone string) error
}

type pointMediaRepository interface {
	Add(ctx context.Context, pointMedia *entity.PointMedia) error
	GetAll(ctx context.Context, pointCode string) ([]*entity.PointMedia, error)
	GetWithMedia(ctx context.Context, pointCode string) ([]*entity.Media, error)
	Get(ctx context.Context, pointCode string, mediaID uuid.UUID) (*entity.PointMedia, error)
	UpdateOrder(ctx context.Context, pointCode string, mediaID uuid.UUID, displayOrder int) error
	Remove(ctx context.Context, pointCode string, mediaID uuid.UUID) error
	RemoveAll(ctx context.Context, pointCode string) error
	Count(ctx context.Context, pointCode string) (int, error)
	Has(ctx context.Context, pointCode string, mediaID uuid.UUID) (bool, error)
}

type Service struct {
	txManager      database.TXManager
	storage        storage
	mediaRepo      mediaRepository
	userAvatarRepo userAvatarRepository
	pointMediaRepo pointMediaRepository
	mediaTTL       time.Duration
}

func NewService(
	txManager database.TXManager,
	storage storage,
	mediaRepo mediaRepository,
	userAvatarRepo userAvatarRepository,
	pointMediaRepo pointMediaRepository,
	mediaTTL time.Duration,
) *Service {
	return &Service{
		txManager:      txManager,
		storage:        storage,
		mediaRepo:      mediaRepo,
		userAvatarRepo: userAvatarRepo,
		pointMediaRepo: pointMediaRepo,
		mediaTTL:       mediaTTL,
	}
}

func (s *Service) GenerateUploadURL(ctx context.Context, filename, contentType, purpose string) (*dto.UploadMediaResponse, error) {
	ses, err := session.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	mediaID := uuid.New()
	storageKey := s.genStorageKey(mediaID, filename, purpose)

	media := &entity.Media{
		ID:               mediaID,
		FileKey:          storageKey,
		BucketName:       s.storage.GetBucketName(ctx),
		OriginalFilename: filename,
		MimeType:         contentType,
		Status:           enum.MediaStatusPending,
		UploadedBy:       ptr.Pointer(ses.Phone()),
	}
	err = s.mediaRepo.Create(ctx, media)
	if err != nil {
		return nil, err
	}

	url, err := s.storage.GeneratePresignedUploadURL(ctx, storageKey, contentType, s.mediaTTL)
	if err != nil {
		return nil, s.mapS3Error(err)
	}
	return &dto.UploadMediaResponse{
		MediaID:   mediaID,
		URL:       url,
		ExpiresAt: time.Now().Add(s.mediaTTL),
	}, nil
}

// GenerateDownloadURL генерирует presigned URL для скачивания по file_key
// Используется когда file_key уже известен (из entity) - БЕЗ запроса к БД
func (s *Service) GenerateDownloadURL(ctx context.Context, fileKey string) (string, error) {
	url, err := s.storage.GeneratePresignedDownloadURL(ctx, fileKey, s.mediaTTL)
	if err != nil {
		return "", s.mapS3Error(err)
	}
	return url, nil
}

func (s *Service) genStorageKey(mediaID uuid.UUID, filename, purpose string) string {
	// 2. Безопасное извлечение расширения
	ext := strings.ToLower(filepath.Ext(filename))

	// 3. Формирование ключа с датой для удобства навигации
	// Результат: avatars/2024/05/20/uuid.jpg
	datePart := time.Now().Format("2006/01/02")
	storageKey := fmt.Sprintf("%s/%s/%s%s", purpose, datePart, mediaID.String(), ext)

	return storageKey
}
