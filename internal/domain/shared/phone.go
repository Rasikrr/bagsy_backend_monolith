package shared

import "strings"

type Phone struct {
	value string
}

func NewPhone(val string) (Phone, error) {
	val = strings.TrimSpace(val)
	if val == "" {
		return Phone{}, ErrInvalidPhone
	}
	if val[0] == '+' {
		val = val[1:]
	}
	return Phone{value: val}, nil
}

func (p Phone) String() string {
	return p.value
}

func (p Phone) IsEmpty() bool {
	return p.value == ""
}
