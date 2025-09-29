// nolint: depguard, gosec
package codegen

import (
	"fmt"
	"math/rand"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateAuthCode() string {
	return fmt.Sprintf("%04d", rng.Intn(10000))
}

func GenerateBagsyID() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = alphabet[rng.Intn(len(alphabet))]
	}
	return string(b)
}
