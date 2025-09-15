// nolint: depguard, gosec
package codegen

import (
	"fmt"
	"math/rand"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateAuthCode() string {
	return fmt.Sprintf("%04d", rng.Intn(10000))
}
