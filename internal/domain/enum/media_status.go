package enum

//go:generate enumer -type=MediaStatus -json -trimprefix MediaStatus -transform=snake -output media_status_enumer.go

// MediaStatus представляет статус обработки медиа-файла
type MediaStatus int8

const (
	// MediaStatusPending - файл загружен в S3, но еще не привязан к сущности
	MediaStatusPending MediaStatus = iota
	// MediaStatusActive - файл активно используется
	MediaStatusActive
	// MediaStatusInactive - файл был заменен или отвязан
	MediaStatusInactive
	// MediaStatusFailed - ошибка при обработке файла
	MediaStatusFailed
)
