package s3

import (
	"bytes"
	"context"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
)

// Config содержит параметры конфигурации для S3 клиента
type Config struct {
	Region          string
	Endpoint        string // Опционально, для совместимых с S3 сервисов (MinIO, LocalStack)
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
}

// Client предоставляет методы для работы с AWS S3
type Client struct {
	s3Client      *s3.Client
	presignClient *s3.PresignClient
	uploader      *manager.Uploader
	downloader    *manager.Downloader
	bucketName    string
}

// NewClient создает новый S3 клиент
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.Region == "" {
		return nil, domainErr.ErrS3EmptyRegion
	}
	if cfg.AccessKeyID == "" {
		return nil, domainErr.ErrS3EmptyAccessKey
	}
	if cfg.SecretAccessKey == "" {
		return nil, domainErr.ErrS3EmptySecretKey
	}
	if cfg.BucketName == "" {
		return nil, domainErr.ErrS3EmptyBucket
	}

	// Настройка AWS конфигурации
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, domainErr.ErrS3ConfigFailed.WithError(err)
	}

	// Создание S3 клиента
	var s3ClientOpts []func(*s3.Options)
	if cfg.Endpoint != "" {
		s3ClientOpts = append(s3ClientOpts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = true // Для совместимости с MinIO и LocalStack
		})
	}

	s3Client := s3.NewFromConfig(awsCfg, s3ClientOpts...)

	// Создание presign клиента для генерации подписанных URL
	presignClient := s3.NewPresignClient(s3Client)

	// Создание uploader и downloader для эффективной работы с большими файлами
	uploader := manager.NewUploader(s3Client, func(u *manager.Uploader) {
		u.PartSize = 10 * 1024 * 1024 // 10MB частями
		u.Concurrency = 5             // 5 параллельных загрузок
	})

	downloader := manager.NewDownloader(s3Client, func(d *manager.Downloader) {
		d.PartSize = 10 * 1024 * 1024 // 10MB частями
		d.Concurrency = 5             // 5 параллельных скачиваний
	})

	return &Client{
		s3Client:      s3Client,
		presignClient: presignClient,
		uploader:      uploader,
		downloader:    downloader,
		bucketName:    cfg.BucketName,
	}, nil
}

// Upload загружает файл в S3
func (c *Client) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	if key == "" {
		return "", domainErr.ErrS3EmptyKey
	}
	if len(data) == 0 {
		return "", domainErr.ErrS3EmptyData
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}

	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	result, err := c.uploader.Upload(ctx, input)
	if err != nil {
		return "", domainErr.ErrS3UploadFailed.WithError(err)
	}

	if result.Location == "" {
		return "", domainErr.ErrS3EmptyLocation
	}

	return result.Location, nil
}

// UploadStream загружает файл из io.Reader в S3
func (c *Client) UploadStream(ctx context.Context, key string, reader io.Reader, contentType string) (string, error) {
	if key == "" {
		return "", domainErr.ErrS3EmptyKey
	}
	if reader == nil {
		return "", domainErr.ErrS3EmptyData
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
		Body:   reader,
	}

	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	result, err := c.uploader.Upload(ctx, input)
	if err != nil {
		return "", domainErr.ErrS3UploadFailed.WithError(err)
	}

	if result.Location == "" {
		return "", domainErr.ErrS3EmptyLocation
	}

	return result.Location, nil
}

// Download скачивает файл из S3
func (c *Client) Download(ctx context.Context, key string) ([]byte, error) {
	if key == "" {
		return nil, domainErr.ErrS3EmptyKey
	}

	buffer := manager.NewWriteAtBuffer([]byte{})

	_, err := c.downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, domainErr.ErrS3DownloadFailed.WithError(err)
	}

	return buffer.Bytes(), nil
}

// DownloadStream скачивает файл из S3 в io.WriterAt
func (c *Client) DownloadStream(ctx context.Context, key string, writer io.WriterAt) (int64, error) {
	if key == "" {
		return 0, domainErr.ErrS3EmptyKey
	}
	if writer == nil {
		return 0, domainErr.NewInvalidInputError("writer cannot be nil", nil)
	}

	n, err := c.downloader.Download(ctx, writer, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return 0, domainErr.ErrS3DownloadFailed.WithError(err)
	}

	return n, nil
}

// Delete удаляет файл из S3
func (c *Client) Delete(ctx context.Context, key string) error {
	if key == "" {
		return domainErr.ErrS3EmptyKey
	}

	_, err := c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return domainErr.ErrS3DeleteFailed.WithError(err)
	}

	return nil
}

// DeleteMultiple удаляет несколько файлов из S3
func (c *Client) DeleteMultiple(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return domainErr.NewInvalidInputError("keys cannot be empty", nil)
	}

	var objectIds []types.ObjectIdentifier
	for _, key := range keys {
		if key != "" {
			objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
		}
	}

	if len(objectIds) == 0 {
		return domainErr.NewInvalidInputError("no valid keys provided", nil)
	}

	_, err := c.s3Client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(c.bucketName),
		Delete: &types.Delete{
			Objects: objectIds,
			Quiet:   aws.Bool(true),
		},
	})
	if err != nil {
		return domainErr.ErrS3DeleteFailed.WithError(err)
	}

	return nil
}

// List возвращает список файлов в S3 бакете с заданным префиксом
func (c *Client) List(ctx context.Context, prefix string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucketName),
	}

	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	var keys []string
	paginator := s3.NewListObjectsV2Paginator(c.s3Client, input)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, domainErr.ErrS3ListFailed.WithError(err)
		}

		for _, obj := range output.Contents {
			if obj.Key != nil {
				keys = append(keys, *obj.Key)
			}
		}
	}

	return keys, nil
}

// ListWithDetails возвращает список файлов с подробной информацией
func (c *Client) ListWithDetails(ctx context.Context, prefix string) ([]ObjectInfo, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucketName),
	}

	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	var objects []ObjectInfo
	paginator := s3.NewListObjectsV2Paginator(c.s3Client, input)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, domainErr.ErrS3ListFailed.WithError(err)
		}

		for _, obj := range output.Contents {
			objects = append(objects, ObjectInfo{
				Key:          aws.ToString(obj.Key),
				Size:         aws.ToInt64(obj.Size),
				LastModified: aws.ToTime(obj.LastModified),
				ETag:         aws.ToString(obj.ETag),
			})
		}
	}

	return objects, nil
}

// Exists проверяет существование файла в S3
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	if key == "" {
		return false, domainErr.ErrS3EmptyKey
	}

	_, err := c.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		// Проверяем, является ли ошибка "NotFound" или "NoSuchKey"
		errMsg := err.Error()
		if strings.Contains(errMsg, "NotFound") || strings.Contains(errMsg, "NoSuchKey") {
			return false, nil
		}
		return false, domainErr.NewInternalError("failed to check object existence", err)
	}

	return true, nil
}

// GetURL возвращает публичный URL файла (работает только если бакет публичный)
func (c *Client) GetURL(key string) string {
	if key == "" {
		return ""
	}
	// Для приватных бакетов нужно использовать presigned URL
	return "https://" + c.bucketName + ".s3.amazonaws.com/" + key
}

// GeneratePresignedUploadURL генерирует подписанный URL для загрузки файла фронтендом
// key - ключ (путь) файла в S3
// contentType - тип содержимого (например, "image/jpeg", "application/pdf")
// expiresIn - время жизни ссылки (рекомендуется 15 минут для загрузки)
func (c *Client) GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, expiresIn time.Duration) (string, error) {
	if key == "" {
		return "", domainErr.ErrS3EmptyKey
	}
	if expiresIn <= 0 {
		return "", domainErr.NewInvalidInputError("expiration time must be positive", nil)
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}

	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	presignedReq, err := c.presignClient.PresignPutObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = expiresIn
	})
	if err != nil {
		return "", domainErr.NewInternalError("failed to generate presigned upload URL", err)
	}

	return presignedReq.URL, nil
}

// GeneratePresignedDownloadURL генерирует подписанный URL для скачивания файла
// key - ключ (путь) файла в S3
// expiresIn - время жизни ссылки (рекомендуется от 15 минут до 7 дней)
func (c *Client) GeneratePresignedDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (string, error) {
	if key == "" {
		return "", domainErr.ErrS3EmptyKey
	}
	if expiresIn <= 0 {
		return "", domainErr.NewInvalidInputError("expiration time must be positive", nil)
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}

	presignedReq, err := c.presignClient.PresignGetObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = expiresIn
	})
	if err != nil {
		return "", domainErr.NewInternalError("failed to generate presigned download URL", err)
	}

	return presignedReq.URL, nil
}
