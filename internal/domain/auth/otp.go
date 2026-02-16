package auth

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
)

const (
	codeLength = 4
)

// OTPCode is an aggregate representing a temporary authentication code.
type OTPCode struct {
	Phone     shared.Phone
	Code      string
	ExpiresAt time.Time
	Attempts  int
}

func NewOTPCode(
	phone shared.Phone,
	duration time.Duration,
) (*OTPCode, error) {
	code, err := generateRandomCode(codeLength)
	if err != nil {
		return nil, err
	}

	return &OTPCode{
		Phone:     phone,
		Code:      code,
		ExpiresAt: time.Now().Add(duration),
		Attempts:  0,
	}, nil
}

func (o *OTPCode) Verify(code string) error {
	if time.Now().After(o.ExpiresAt) {
		return ErrOTPExpired
	}
	if o.Code != code {
		o.Attempts++
		return ErrOTPInvalid
	}
	return nil
}

func (o *OTPCode) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

func generateRandomCode(length int) (string, error) {
	const charset = "0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}
