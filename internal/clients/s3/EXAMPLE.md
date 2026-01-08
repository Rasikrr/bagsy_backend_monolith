# Пример создания API для загрузки файлов

## Backend Handler (Go)

```go
package files

import (
    "fmt"
    "net/http"
    "time"

    "github.com/google/uuid"
    s3Client "github.com/Rasikrr/bagsy_backend_monolith/internal/clients/s3"
    domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
    "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
    "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/httputil"
)

type S3Service interface {
    GeneratePresignedUploadURL(ctx context.Context, key, contentType string, expiresIn time.Duration) (string, error)
    GeneratePresignedDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (string, error)
}

type Controller struct {
    s3Service S3Service
}

func New(s3Svc S3Service) *Controller {
    return &Controller{s3Service: s3Svc}
}

func (c *Controller) Init(router *chi.Mux) {
    router.Route("/api/v1/files", func(r chi.Router) {
        r.Post("/upload-url", c.generateUploadURL)
        r.Get("/download-url/{key}", c.generateDownloadURL)
    })
}

// Request/Response models
//go:generate easyjson -all models.go

type generateUploadURLRequest struct {
    FileName    string `json:"fileName" validate:"required"`
    ContentType string `json:"contentType" validate:"required"`
    Folder      string `json:"folder"`  // Опционально: avatars, documents, etc.
}

func (r *generateUploadURLRequest) validate() error {
    return GetValidator().Struct(r)
}

type uploadURLResponse struct {
    UploadURL string `json:"uploadUrl"`
    FileKey   string `json:"fileKey"`   // Путь файла в S3 для последующих операций
    ExpiresIn int    `json:"expiresIn"` // Секунды до истечения
}

type downloadURLResponse struct {
    DownloadURL string `json:"downloadUrl"`
    ExpiresIn   int    `json:"expiresIn"`
}

// @Summary Генерация URL для загрузки файла
// @Description Создает подписанный URL для прямой загрузки файла в S3 с фронтенда
// @Tags files
// @Accept json
// @Produce json
// @Param request body generateUploadURLRequest true "Параметры файла"
// @Success 200 {object} api.SuccessResponse{data=uploadURLResponse}
// @Failure 400,500 {object} api.ErrorResponse
// @Router /api/v1/files/upload-url [post]
func (c *Controller) generateUploadURL(w http.ResponseWriter, r *http.Request) {
    var req generateUploadURLRequest
    ctx := r.Context()

    if err := httputil.GetData(r, &req); err != nil {
        errors.HandleError(ctx, w, err)
        return
    }

    if err := req.validate(); err != nil {
        errors.HandleError(ctx, w, err)
        return
    }

    // Генерируем уникальный ключ для файла
    fileID := uuid.New().String()
    folder := "uploads"
    if req.Folder != "" {
        folder = req.Folder
    }
    fileKey := fmt.Sprintf("%s/%s-%s", folder, fileID, req.FileName)

    // Генерируем presigned URL (15 минут)
    expiresIn := 15 * time.Minute
    uploadURL, err := c.s3Service.GeneratePresignedUploadURL(ctx, fileKey, req.ContentType, expiresIn)
    if err != nil {
        errors.HandleError(ctx, w, err)
        return
    }

    httputil.SendData(ctx, w, uploadURLResponse{
        UploadURL: uploadURL,
        FileKey:   fileKey,
        ExpiresIn: int(expiresIn.Seconds()),
    }, http.StatusOK)
}

// @Summary Генерация URL для скачивания файла
// @Description Создает подписанный URL для безопасного скачивания файла из S3
// @Tags files
// @Produce json
// @Param key path string true "Ключ файла в S3"
// @Success 200 {object} api.SuccessResponse{data=downloadURLResponse}
// @Failure 400,404,500 {object} api.ErrorResponse
// @Router /api/v1/files/download-url/{key} [get]
func (c *Controller) generateDownloadURL(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    key := chi.URLParam(r, "key")

    if key == "" {
        errors.HandleError(ctx, w, domainErr.ErrS3EmptyKey)
        return
    }

    // Генерируем presigned URL (1 час)
    expiresIn := 1 * time.Hour
    downloadURL, err := c.s3Service.GeneratePresignedDownloadURL(ctx, key, expiresIn)
    if err != nil {
        errors.HandleError(ctx, w, err)
        return
    }

    httputil.SendData(ctx, w, downloadURLResponse{
        DownloadURL: downloadURL,
        ExpiresIn:   int(expiresIn.Seconds()),
    }, http.StatusOK)
}
```

## Frontend (React/TypeScript)

### Хук для загрузки файлов

```typescript
// hooks/useFileUpload.ts
import { useState } from 'react';

interface UploadURLResponse {
  uploadUrl: string;
  fileKey: string;
  expiresIn: number;
}

interface UseFileUploadResult {
  upload: (file: File, folder?: string) => Promise<string>;
  uploading: boolean;
  progress: number;
  error: string | null;
}

export const useFileUpload = (): UseFileUploadResult => {
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);

  const upload = async (file: File, folder: string = 'uploads'): Promise<string> => {
    setUploading(true);
    setProgress(0);
    setError(null);

    try {
      // 1. Получаем presigned URL от бэкенда
      const urlResponse = await fetch('/api/v1/files/upload-url', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          fileName: file.name,
          contentType: file.type,
          folder: folder,
        }),
      });

      if (!urlResponse.ok) {
        throw new Error('Failed to get upload URL');
      }

      const { uploadUrl, fileKey }: UploadURLResponse = await urlResponse.json();

      // 2. Загружаем файл напрямую в S3
      const xhr = new XMLHttpRequest();

      // Отслеживаем прогресс
      xhr.upload.addEventListener('progress', (e) => {
        if (e.lengthComputable) {
          setProgress(Math.round((e.loaded / e.total) * 100));
        }
      });

      // Промисифицируем XMLHttpRequest для await
      const uploadPromise = new Promise<void>((resolve, reject) => {
        xhr.addEventListener('load', () => {
          if (xhr.status === 200) {
            resolve();
          } else {
            reject(new Error(`Upload failed with status ${xhr.status}`));
          }
        });

        xhr.addEventListener('error', () => {
          reject(new Error('Upload failed'));
        });
      });

      xhr.open('PUT', uploadUrl);
      xhr.setRequestHeader('Content-Type', file.type);
      xhr.send(file);

      await uploadPromise;

      setProgress(100);
      setUploading(false);

      return fileKey;
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Upload failed';
      setError(message);
      setUploading(false);
      throw err;
    }
  };

  return { upload, uploading, progress, error };
};
```

### Компонент загрузки файла

```typescript
// components/FileUploader.tsx
import React, { useRef } from 'react';
import { useFileUpload } from '../hooks/useFileUpload';

interface FileUploaderProps {
  folder?: string;
  onUploadComplete?: (fileKey: string) => void;
  accept?: string;
}

export const FileUploader: React.FC<FileUploaderProps> = ({
  folder = 'uploads',
  onUploadComplete,
  accept,
}) => {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const { upload, uploading, progress, error } = useFileUpload();

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    try {
      const fileKey = await upload(file, folder);
      console.log('File uploaded successfully:', fileKey);
      onUploadComplete?.(fileKey);
    } catch (err) {
      console.error('Upload error:', err);
    }
  };

  return (
    <div className="file-uploader">
      <input
        ref={fileInputRef}
        type="file"
        onChange={handleFileSelect}
        accept={accept}
        disabled={uploading}
        style={{ display: 'none' }}
      />

      <button
        onClick={() => fileInputRef.current?.click()}
        disabled={uploading}
      >
        {uploading ? `Uploading... ${progress}%` : 'Choose File'}
      </button>

      {error && <div className="error">{error}</div>}
      {uploading && <div className="progress-bar" style={{ width: `${progress}%` }} />}
    </div>
  );
};
```

### Использование компонента

```typescript
// App.tsx
import React from 'react';
import { FileUploader } from './components/FileUploader';

export const App: React.FC = () => {
  const handleUploadComplete = (fileKey: string) => {
    console.log('File uploaded with key:', fileKey);
    // Сохраняем fileKey в базу данных или используем для других операций
  };

  return (
    <div>
      <h1>Upload Avatar</h1>
      <FileUploader
        folder="avatars"
        accept="image/*"
        onUploadComplete={handleUploadComplete}
      />
    </div>
  );
};
```

## Сохранение метаданных файла

После успешной загрузки фронтенд отправляет fileKey на бэкенд для сохранения метаданных:

```typescript
// После загрузки файла
const handleUploadComplete = async (fileKey: string) => {
  // Сохраняем метаданные в базе данных через API
  await fetch('/api/v1/users/avatar', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      avatarKey: fileKey,
    }),
  });
};
```

Backend обновляет базу данных:

```go
type updateAvatarRequest struct {
    AvatarKey string `json:"avatarKey" validate:"required"`
}

func (c *Controller) updateAvatar(w http.ResponseWriter, r *http.Request) {
    var req updateAvatarRequest
    ctx := r.Context()

    // ... валидация и получение пользователя из контекста

    // Обновляем avatar_key в базе данных
    err := c.userService.UpdateAvatar(ctx, userID, req.AvatarKey)
    if err != nil {
        errors.HandleError(ctx, w, err)
        return
    }

    httputil.SendSuccess(ctx, w, "Avatar updated successfully")
}
```

## Преимущества этого подхода

1. **Экономия трафика**: Файлы загружаются напрямую в S3, минуя бэкенд
2. **Масштабируемость**: Бэкенд не тратит ресурсы на передачу файлов
3. **Производительность**: Загрузка быстрее за счёт прямого соединения с S3
4. **Безопасность**: URL подписан и имеет ограниченный срок действия
5. **Простота**: Фронтенд работает с обычным HTTP PUT запросом