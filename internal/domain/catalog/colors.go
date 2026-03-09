package catalog

type Color string

const (
	ColorUnknown Color = "unknown"
	ColorBlue    Color = "blue"
	ColorGreen   Color = "green"
	ColorRed     Color = "red"
	ColorYellow  Color = "yellow"
	ColorPurple  Color = "purple"
	ColorOrange  Color = "orange"
	ColorGray    Color = "gray"
)

func (c Color) String() string {
	return string(c)
}

func (c Color) IsValid() bool {
	switch c {
	case ColorBlue, ColorGreen, ColorRed, ColorYellow, ColorPurple, ColorOrange, ColorGray:
		return true
	default:
		return false
	}
}

func ParseColor(s string) (Color, error) {
	c := Color(s)
	if !c.IsValid() {
		return ColorUnknown, ErrServiceInvalidColor
	}
	return c, nil
}
