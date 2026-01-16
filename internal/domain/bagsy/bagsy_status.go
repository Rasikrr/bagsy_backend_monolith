package bagsy

//go:generate enumer -type=Status -json -trimprefix Status -transform=snake -output bagsy_status_enumer.go

type Status uint8

const (
	StatusPending Status = iota
	StatusCreated
	StatusActive
	StatusCompleted
	StatusCanceled
)
