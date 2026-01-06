package slug

import (
	"sync"

	"github.com/gosimple/slug"
)

var (
	once sync.Once
)

// nolint: reassign
func setupSlug() {
	once.Do(func() {
		slug.CustomSub = map[string]string{
			" ":  "_",
			",":  "_",
			".":  "_",
			"-":  "_",
			"/":  "_",
			"\"": "_",
		}
		slug.Lowercase = true
	})
}

func Generate(s string) string {
	setupSlug()
	result := slug.Make(s)
	return result
}
