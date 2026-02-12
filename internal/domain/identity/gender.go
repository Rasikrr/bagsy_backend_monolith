package identity

type Gender string

const (
	GenderMale        Gender = "male"
	GenderFemale      Gender = "female"
	GenderUnspecified Gender = "unspecified"
)

func (g Gender) IsValid() bool {
	switch g {
	case GenderMale, GenderFemale, GenderUnspecified:
		return true
	}
	return false
}

func (g Gender) String() string {
	return string(g)
}

func ParseGender(s string) (Gender, error) {
	g := Gender(s)
	if !g.IsValid() {
		return GenderUnspecified, ErrInvalidGender
	}
	return g, nil
}
