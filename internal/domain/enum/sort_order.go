package enum

//go:generate enumer -type=SortOrder -json -trimprefix SortOrder -transform=snake -output sort_order_enumer.go

type SortOrder uint8

const (
	SortOrderDesc SortOrder = iota
	SortOrderAsc
)
