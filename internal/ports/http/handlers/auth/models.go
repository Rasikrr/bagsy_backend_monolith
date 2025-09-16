package auth

//go:generate easyjson -all models.go
type sendCodeRequest struct {
	Phone string `json:"phone"`
}
