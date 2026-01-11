package enum

//go:generate enumer -type=MediaStatus -json -trimprefix MediaStatus -transform=snake_case

// MediaStatus представляет статус обработки медиа-файла
type MediaStatus int8

const (
	// MediaStatusPending - файл загружен в S3, но еще не привязан к сущности
	MediaStatusPending MediaStatus = iota
	// MediaStatusActive - файл активно используется
	MediaStatusActive
	// MediaStatusProcessing - файл в процессе обработки (например, генерация thumbnails)
	MediaStatusProcessing
	// MediaStatusFailed - ошибка при обработке файла
	MediaStatusFailed
)
