package enum

//go:generate enumer -type=BagsyStatus -json -trimprefix BagsyStatus -transform=snake -output bagsy_status_enumer.go

type BagsyStatus uint8

const (
	BagsyStatusPending BagsyStatus = iota
	BagsyStatusCreated
	BagsyStatusActive
	BagsyStatusCompleted
	BagsyStatusCanceled
)
