package slug

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlug(t *testing.T) {
	type testCase struct {
		Input string
		Want  string
	}

	testCases := []testCase{
		{
			Input: "a",
			Want:  "a",
		},
		{
			Input: "Ернар Ханапин",
			Want:  "ernar_khanapin",
		},
		{
			Input: " Дмитрий Каиргелдин ",
			Want:  "dmitrii_kairgeldin",
		},
		{
			Input: " Rassul Turtulov. ",
			Want:  "rassul_turtulov",
		},
	}
	for _, tc := range testCases {
		assert.Equal(t, tc.Want, Generate(tc.Input))
	}
}
