package media

// Status определяет текущее состояние файла в процессе загрузки (Direct-to-S3)
type Status string

const (
	StatusPending  Status = "pending"
	StatusUploaded Status = "uploaded"
	StatusFailed   Status = "failed"
)
