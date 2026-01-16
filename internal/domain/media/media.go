package media

import (
	"time"

	"github.com/google/uuid"
)

// Media представляет медиа-файл (изображение) в системе
type Media struct {
	ID               uuid.UUID
	FileKey          string // Путь в S3: "YYYY/MM/{uuid}.{ext}"
	BucketName       string
	OriginalFilename string
	MimeType         string // image/jpeg, image/png, image/webp
	SizeBytes        int64
	Width            *int // Для изображений
	Height           *int // Для изображений
	Status           Status
	UploadedBy       *string // Кто загрузил файл
	CreatedAt        time.Time
	UpdatedAt        *time.Time
	DeletedAt        *time.Time
}

// UserMedia представляет связь пользователя с его аватаром
type UserMedia struct {
	UserPhone string
	MediaID   uuid.UUID
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// PointMedia представляет связь точки с её фотографиями
type PointMedia struct {
	ID           uuid.UUID
	PointCode    string
	MediaID      uuid.UUID
	DisplayOrder int
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time
}
