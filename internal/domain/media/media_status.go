package media

//go:generate enumer -type=Status -json -trimprefix Status -transform=snake -output media_status_enumer.go

// Status MediaStatus представляет статус обработки медиа-файла
type Status int8

const (
	// StatusPending MediaStatusPending - файл загружен в S3, но еще не привязан к сущности
	StatusPending Status = iota
	// StatusActive MediaStatusActive - файл активно используется
	StatusActive
	// StatusInactive MediaStatusInactive - файл был заменен или отвязан
	StatusInactive
	// StatusFailed MediaStatusFailed - ошибка при обработке файла
	StatusFailed
)
