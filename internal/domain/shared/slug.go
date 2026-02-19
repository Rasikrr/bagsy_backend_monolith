package shared

import "github.com/Rasikrr/bagsy_backend_monolith/internal/util/slug"

type Slug struct {
	value string
}

func NewSlug(value string) (Slug, error) {
	if value == "" {
		return Slug{}, ErrEmptySlug
	}
	return Slug{
		value: slug.Generate(value),
	}, nil
}

func ParseSlug(slug string) (Slug, error) {
	if slug == "" {
		return Slug{}, ErrEmptySlug
	}
	return Slug{value: slug}, nil
}

func (s Slug) String() string {
	return s.value
}

func (s Slug) IsEmpty() bool {
	return s.String() == ""
}

func (s Slug) IsEqual(s2 Slug) bool {
	return s.value == s2.value
}
