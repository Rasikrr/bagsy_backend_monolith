package auth

//go:generate enumer -type=TokenPurpose -json -trimprefix TokenPurpose -transform=snake -output token_purpose_enumer.go

type TokenPurpose int

const (
	TokenPurposeRegister TokenPurpose = iota
	TokenPurposePasswordChange
)
