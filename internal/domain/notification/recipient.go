package notification

//go:generate enumer -type=RecipientType -json -trimprefix RecipientType -transform=snake -output recipient_type_enumer.go
type RecipientType int8

const (
	RecipientTypeClient RecipientType = iota
	RecipientTypeMaster
)

func AllRecipients() []RecipientType {
	return []RecipientType{RecipientTypeClient, RecipientTypeMaster}
}
