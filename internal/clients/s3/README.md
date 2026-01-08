# S3 Client

Клиент для работы с AWS S3 (и совместимыми сервисами) на основе AWS SDK Go v2.

## Возможности

- Загрузка файлов (Upload, UploadStream)
- Скачивание файлов (Download, DownloadStream)
- Удаление файлов (Delete, DeleteMultiple)
- Получение списка файлов (List, ListWithDetails)
- Проверка существования файла (Exists)
- Получение публичного URL (GetURL)
- **Генерация Presigned URLs** для прямой загрузки/скачивания фронтендом
- Поддержка multipart загрузки для больших файлов
- Совместимость с MinIO и LocalStack

## Конфигурация

Добавьте следующие переменные окружения в `.env`:

```env
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_S3_BUCKET_NAME=your_bucket_name
AWS_S3_ENDPOINT=         # Опционально, для MinIO или LocalStack
```

## Использование

### Создание клиента

```go
import (
    "context"
    s3Client "github.com/Rasikrr/bagsy_backend_monolith/internal/clients/s3"
)

ctx := context.Background()

client, err := s3Client.NewClient(ctx, s3Client.Config{
    Region:          "us-east-1",
    AccessKeyID:     "your_access_key",
    SecretAccessKey: "your_secret_key",
    BucketName:      "your_bucket_name",
    Endpoint:        "", // Оставьте пустым для AWS S3, укажите для MinIO/LocalStack
})
if err != nil {
    log.Fatal(err)
}
```

### Загрузка файла

```go
// Загрузка из байтов
data := []byte("Hello, S3!")
location, err := client.Upload(ctx, "files/hello.txt", data, "text/plain")
if err != nil {
    log.Fatal(err)
}
fmt.Println("Uploaded to:", location)

// Загрузка из Reader
file, _ := os.Open("document.pdf")
defer file.Close()
location, err = client.UploadStream(ctx, "documents/doc.pdf", file, "application/pdf")
```

### Скачивание файла

```go
// Скачивание в память
data, err := client.Download(ctx, "files/hello.txt")
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(data))

// Скачивание в файл
file, _ := os.Create("downloaded.pdf")
defer file.Close()
n, err := client.DownloadStream(ctx, "documents/doc.pdf", file)
fmt.Printf("Downloaded %d bytes\n", n)
```

### Удаление файлов

```go
// Удаление одного файла
err := client.Delete(ctx, "files/hello.txt")

// Удаление нескольких файлов
err = client.DeleteMultiple(ctx, []string{
    "files/file1.txt",
    "files/file2.txt",
    "files/file3.txt",
})
```

### Получение списка файлов

```go
// Простой список ключей
keys, err := client.List(ctx, "files/")
for _, key := range keys {
    fmt.Println(key)
}

// Список с подробной информацией
objects, err := client.ListWithDetails(ctx, "files/")
for _, obj := range objects {
    fmt.Printf("%s - %d bytes - %v\n", obj.Key, obj.Size, obj.LastModified)
}
```

### Проверка существования файла

```go
exists, err := client.Exists(ctx, "files/hello.txt")
if err != nil {
    log.Fatal(err)
}
if exists {
    fmt.Println("File exists!")
}
```

### Получение URL

```go
url := client.GetURL("files/hello.txt")
fmt.Println("Public URL:", url)
// Для приватных бакетов используйте presigned URL
```

### Генерация Presigned URLs для фронтенда

#### Генерация URL для загрузки файла

Бэкенд генерирует подписанный URL, который фронтенд использует для прямой загрузки в S3:

```go
// Генерируем URL для загрузки файла (15 минут)
uploadURL, err := client.GeneratePresignedUploadURL(
    ctx,
    "uploads/avatar-123.jpg",      // Путь в S3
    "image/jpeg",                   // Content-Type
    15 * time.Minute,               // Время жизни ссылки
)
if err != nil {
    log.Fatal(err)
}

// Возвращаем URL фронтенду
fmt.Println("Upload URL:", uploadURL)
```

**Использование на фронтенде (JavaScript/TypeScript):**

```javascript
// Получаем URL от бэкенда
const response = await fetch('/api/upload-url', {
    method: 'POST',
    body: JSON.stringify({ fileName: 'avatar.jpg', contentType: 'image/jpeg' })
});
const { uploadUrl } = await response.json();

// Загружаем файл напрямую в S3
const file = document.getElementById('fileInput').files[0];
const uploadResponse = await fetch(uploadUrl, {
    method: 'PUT',
    headers: {
        'Content-Type': 'image/jpeg'
    },
    body: file
});

if (uploadResponse.ok) {
    console.log('File uploaded successfully!');
}
```

#### Генерация URL для скачивания файла

Для безопасного скачивания файлов из приватного бакета:

```go
// Генерируем URL для скачивания (1 час)
downloadURL, err := client.GeneratePresignedDownloadURL(
    ctx,
    "documents/report-456.pdf",
    1 * time.Hour,
)
if err != nil {
    log.Fatal(err)
}

// Возвращаем URL фронтенду
fmt.Println("Download URL:", downloadURL)
```

**Использование на фронтенде:**

```javascript
// Получаем URL от бэкенда
const response = await fetch('/api/download-url/report-456.pdf');
const { downloadUrl } = await response.json();

// Открываем файл или скачиваем
window.open(downloadUrl, '_blank');
// или
window.location.href = downloadUrl;
```

#### Рекомендации по времени жизни ссылок

- **Загрузка файлов**: 15-30 минут (пользователь выбирает файл и загружает)
- **Скачивание файлов**: 1-24 часа (в зависимости от use case)
- **Временные ссылки**: 5-15 минут (для превью, временного доступа)
- **Долгосрочный доступ**: до 7 дней (максимум для AWS S3)

```go
const (
    UploadExpiration   = 15 * time.Minute    // Для загрузки
    DownloadExpiration = 1 * time.Hour       // Для скачивания
    PreviewExpiration  = 5 * time.Minute     // Для превью
    ShareExpiration    = 24 * time.Hour      // Для публичных ссылок
)
```

## Обработка ошибок

Клиент возвращает доменные ошибки из `internal/domain/errors`:

```go
data, err := client.Download(ctx, "nonexistent.txt")
if err != nil {
    if domainErr.IsNotFound(err) {
        fmt.Println("File not found")
    } else {
        fmt.Println("Error:", err)
    }
}
```

Доступные ошибки валидации:
- `ErrS3EmptyRegion`
- `ErrS3EmptyAccessKey`
- `ErrS3EmptySecretKey`
- `ErrS3EmptyBucket`
- `ErrS3EmptyKey`
- `ErrS3EmptyData`

Ошибки API:
- `ErrS3ConfigFailed`
- `ErrS3UploadFailed`
- `ErrS3DownloadFailed`
- `ErrS3DeleteFailed`
- `ErrS3ListFailed`
- `ErrS3EmptyLocation`

## Производительность

Клиент использует AWS SDK Manager для эффективной работы с большими файлами:
- Multipart загрузка/скачивание по 10MB частями
- 5 параллельных горутин для одновременной обработки частей
- Автоматическое управление буферами

Для настройки параметров производительности отредактируйте конфигурацию uploader/downloader в `client.go`.

## Тестирование с LocalStack

Для локальной разработки можно использовать LocalStack:

```bash
docker run -d -p 4566:4566 localstack/localstack
```

Конфигурация:
```go
client, err := s3Client.NewClient(ctx, s3Client.Config{
    Region:          "us-east-1",
    AccessKeyID:     "test",
    SecretAccessKey: "test",
    BucketName:      "test-bucket",
    Endpoint:        "http://localhost:4566",
})
```

## Зависимости

- `github.com/aws/aws-sdk-go-v2/aws`
- `github.com/aws/aws-sdk-go-v2/config`
- `github.com/aws/aws-sdk-go-v2/credentials`
- `github.com/aws/aws-sdk-go-v2/feature/s3/manager`
- `github.com/aws/aws-sdk-go-v2/service/s3`
