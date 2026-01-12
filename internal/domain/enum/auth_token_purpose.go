package enum

//go:generate enumer -type=AuthTokenPurpose -json -trimprefix AuthTokenPurpose -transform=snake -output auth_token_purpose_enumer.go

type AuthTokenPurpose int

const (
	AuthTokenPurposeRegister AuthTokenPurpose = iota
	AuthTokenPurposePasswordChange
)
