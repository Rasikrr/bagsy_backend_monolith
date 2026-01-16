package media

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/actor"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
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
	Create(ctx context.Context, media *media.Media) error
	GetByID(ctx context.Context, mediaID uuid.UUID) (*media.Media, error)
	GetByIDs(ctx context.Context, ids ...uuid.UUID) ([]*media.Media, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status media.Status) error
	UpdateStatuses(ctx context.Context, ids []uuid.UUID, status media.Status) error
	SoftDeleteByID(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	txManager database.TXManager
	storage   storage
	mediaRepo mediaRepository
	mediaTTL  time.Duration
}

func NewService(
	txManager database.TXManager,
	storage storage,
	mediaRepo mediaRepository,
	mediaTTL time.Duration,
) *Service {
	return &Service{
		txManager: txManager,
		storage:   storage,
		mediaRepo: mediaRepo,
		mediaTTL:  mediaTTL,
	}
}

func (s *Service) GenerateUploadURL(ctx context.Context, filename, contentType, purpose string) (*media.UploadedMedia, error) {
	ses, err := actor.GetActor(ctx)
	if err != nil {
		return nil, err
	}
	mediaID := uuid.New()
	storageKey := s.genStorageKey(mediaID, filename, purpose)

	med := &media.Media{
		ID:               mediaID,
		FileKey:          storageKey,
		BucketName:       s.storage.GetBucketName(ctx),
		OriginalFilename: filename,
		MimeType:         contentType,
		Status:           media.StatusPending,
		UploadedBy:       ptr.Pointer(ses.Phone()),
	}
	err = s.mediaRepo.Create(ctx, med)
	if err != nil {
		return nil, err
	}

	url, err := s.storage.GeneratePresignedUploadURL(ctx, storageKey, contentType, s.mediaTTL)
	if err != nil {
		return nil, s.mapS3Error(err)
	}
	return &media.UploadedMedia{
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

func (s *Service) GenerateDownloadURLs(ctx context.Context, fileKeys []string) ([]string, error) {
	if len(fileKeys) == 0 {
		return []string{}, nil
	}
	out := make([]string, len(fileKeys))
	for i, key := range fileKeys {
		url, err := s.storage.GeneratePresignedDownloadURL(ctx, key, s.mediaTTL)
		if err != nil {
			return nil, s.mapS3Error(err)
		}
		out[i] = url
	}
	return out, nil
}

func (s *Service) GetByID(ctx context.Context, mediaID uuid.UUID) (*media.Media, error) {
	med, err := s.mediaRepo.GetByID(ctx, mediaID)
	if err != nil {
		return nil, err
	}
	return med, nil
}

func (s *Service) GetByIDs(ctx context.Context, ids ...uuid.UUID) ([]*media.Media, error) {
	med, err := s.mediaRepo.GetByIDs(ctx, ids...)
	if err != nil {
		return nil, err
	}
	return med, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id uuid.UUID, status media.Status) error {
	err := s.mediaRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateStatuses(ctx context.Context, ids []uuid.UUID, status media.Status) error {
	err := s.mediaRepo.UpdateStatuses(ctx, ids, status)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) SoftDeleteByID(ctx context.Context, id uuid.UUID) error {
	err := s.mediaRepo.SoftDeleteByID(ctx, id)
	return err
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
